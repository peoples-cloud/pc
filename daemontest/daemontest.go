package daemontest

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/peoples-cloud/pc/cryptotest"
	"github.com/peoples-cloud/pc/dockertest"
	"github.com/peoples-cloud/pc/ipfstest"
	"github.com/peoples-cloud/pc/rpctest"
	"github.com/peoples-cloud/pc/tartest"
	"github.com/peoples-cloud/pc/util"

	"github.com/BurntSushi/toml"
	"github.com/thoj/go-ircevent"
)

// for rpc between client & server https://golang.org/pkg/net/rpc/

type Listener int
type Config struct {
	Server          string
	Channel         string
	Nick            string
	DeployDirectory string
	RPCPort         string
}

var conn *irc.Connection
var config Config

const IPFS_LENGTH = 46
const KEY_LENGTH = 44
const MAX_ROLL = 16777216
const MIN_ROLL = 0

var isDeploying = false
var myRoll int64 = -1
var highestRoll int64 = -1
var deployTimeout = time.Second * 5
var deployRequest rpctest.Program

// TODO: keep track of deployed programs in bookkeeping
var connectedSwarms = make(map[string]string)

func ReadConfig(configfile string) Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	config.DeployDirectory, err = filepath.Abs(config.DeployDirectory)
	_, err = os.Stat(config.DeployDirectory)
	if err != nil {
		log.Print("config: DeployDirectory doesn't exist; creating it")
		_ = util.MakeDir(config.DeployDirectory, "")
	}
	// if default nick hasn't been changed, assign a random one
	// TODO: handle nick collisions
	if config.Nick == "default" {
		config.Nick = "pc-" + cryptotest.GenerateRandomString(7)
		log.Printf("default nick detected, chose %s as nick\n", config.Nick)
	}

	connectedSwarms[config.Channel] = ""
	return config
}

func joinConnectedSwarms() {
	for swarm, password := range connectedSwarms {
		conn.Join(fmt.Sprintf("#%s %s", swarm, password))
		log.Printf("joined #%s\n", swarm)
	}
}

func updateBookkeeping(info *rpctest.RPCInfo) {
	log.Println("updating bookkeeping")
	password := info.Password
	if len(info.Password) == 0 {
		password = "<none>"
	}
	connectedSwarms[info.Swarm] = password
	data, err := json.Marshal(connectedSwarms)
	if err != nil {
		log.Printf("update bookkeeping: %v\n", err)
		panic(err)
	}
	err = ioutil.WriteFile(config.DeployDirectory+"/.pcinfo", data, os.FileMode(0777))
	if err != nil {
		log.Printf("update bookkeeping: %v\n", err)
		panic(err)
	}
	log.Println("successfully updated bookkeeping")
}

func readBookkeeping() {
	pcfile := config.DeployDirectory + "/.pcinfo"
	data, err := ioutil.ReadFile(pcfile)
	if err != nil {
		log.Printf("couldn't find %s; creating it\n", pcfile)
		err = createBookkeeping(pcfile)
		if err != nil {
			panic(err)
			log.Printf("couldn't create %s\n", pcfile)
			return
		}
		log.Printf("created %s\n", pcfile)
		return
	}
	json.Unmarshal(data, &connectedSwarms)
}

func createBookkeeping(path string) error {
	err := ioutil.WriteFile(path, make([]byte, 0), os.FileMode(0777))
	return err
}

func (l *Listener) Deploy(info *rpctest.RPCInfo, reply *rpctest.Message) error {
	_, ok := connectedSwarms[info.Swarm]
	if !ok {
		reply.Msg = fmt.Sprintf("Not connected to any swarm named %s.", info.Swarm)
		return nil
	}

	isDeploying = true
	log.Printf("RPC received: %s: %s %s\n", info.Swarm, info.Path, info.Language)
	log.Println("creating dockerfile")
	// create dockerfile
	dockertest.CreateDockerfile(info.Language, info.Path)
	log.Printf("wrote dockerfile for %s\n", info.Language)
	// tar destination
	log.Println("creating tarball")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("deploy: %v\n", err)
	}
	tarball := fmt.Sprintf("%s/%s", dir, "pc-docker.tar.gz")
	tartest.PackTar(info.Path, tarball)
	// encrypt tar
	log.Println("encrypting tarball")
	key, tarball := cryptotest.Encrypt(tarball)
	log.Println("dest", tarball)
	log.Printf("key: %s\n", key)
	// upload to ipfs
	log.Println("uploading to ipfs")
	hash := ipfstest.IPFSAdd(tarball)
	log.Printf("hash: %s\n", hash)
	// deploy over irc with hash & key
	channel := fmt.Sprintf("#%s", info.Swarm)
	conn.Privmsg(channel, fmt.Sprintf("[pc-deploy] %s %s", hash, key))
	// reply to user with hash & key
	reply.Msg = fmt.Sprintf("hash: %s\tkey: %s", hash, key)

	go func() {
		time.Sleep(deployTimeout)
		log.Printf("highest roll was %d\n", highestRoll)
		if highestRoll > 0 {
			conn.Privmsg(channel, fmt.Sprintf("[pc-host] %d", highestRoll))
		} else {
			log.Printf("could not deploy %s", hash)
		}
		isDeploying = false
		highestRoll = -1 // reset roll counter
	}()
	return nil
}

func (l *Listener) Create(info *rpctest.RPCInfo, reply *rpctest.Message) error {
	if len(info.Swarm) == 0 {
		info.Swarm = cryptotest.GenerateRandomString(8)
	}
	log.Printf("RPC received: join #%s pass: %s\n", info.Swarm, info.Password)
	conn.Join(fmt.Sprintf("#%s %s", info.Swarm, info.Password))
	reply.Msg = info.Swarm
	updateBookkeeping(info)
	return nil
}

func stopContainer(program rpctest.Program) {
	dockertest.StopContainer(program.Hash)
	conn.Privmsg(fmt.Sprintf("#%s", program.Swarm), fmt.Sprintf("[pc-stopped] %s", program.Hash))
}

func deployFromNetwork(program rpctest.Program) {
	log.Printf("hash: %s\nkey: %s\n", program.Hash, program.Key)
	// create folder
	deployPath := util.MakeDir(config.DeployDirectory, program.Hash+"-deploy")
	// get from ipfs
	log.Println("downloading program from ipfs")
	ipfstest.IPFSGet(program.Hash, deployPath)
	tarball := fmt.Sprintf("%s/%s", deployPath, program.Hash)
	// decrypt
	log.Println("decrypting tar")
	cryptotest.Decrypt(tarball, program.Key, tarball)
	// untar
	log.Println("unpacking tar")
	tartest.UnpackTar(tarball, deployPath)
	// build docker image
	log.Println("building docker image")
	dockertest.BuildImage(deployPath, program.Hash)
	log.Printf("built %s-image\n", program.Hash)
	// start container
	log.Println("starting container")
	err := dockertest.RunContainer(program.Hash)
	if err != nil {
		return
	}
	log.Printf("started %s-container\n", program.Hash)
	conn.Privmsg(fmt.Sprintf("#%s", program.Swarm), fmt.Sprintf("[pc-ack] %s", program.Hash))
}

func (l *Listener) Test(info *rpctest.RPCInfo, reply *rpctest.Message) error {
	// deployFromNetwork(rpctest.Program{Hash: info.Hash, Key: info.Key})
	return nil
}

func (l *Listener) Join(info *rpctest.RPCInfo, reply *rpctest.Message) error {
	log.Printf("RPC received: join #%s pass: %s\n", info.Swarm, info.Password)
	conn.Join(fmt.Sprintf("#%s %s", info.Swarm, info.Password))
	updateBookkeeping(info)
	return nil
}

func (l *Listener) Leave(info *rpctest.RPCInfo, reply *rpctest.Message) error {
	log.Printf("RPC received: part #%s\n", info.Swarm)
	conn.Part(fmt.Sprintf("#%s", info.Swarm))
	updateBookkeeping(info)
	return nil
}

func (l *Listener) Stop(info *rpctest.RPCInfo, reply *rpctest.Message) error {
	log.Printf("RPC received: stop %s\n", info.Hash)
	conn.Privmsg(fmt.Sprintf("#%s", info.Swarm), fmt.Sprintf("[pc-stop] %s", info.Hash))
	return nil
}

func (l *Listener) List(line string, reply *rpctest.Message) error {
	log.Printf("RPC received: %s\n", string(line))
	return nil
}

func rollForDeploy() {
	nBig, err := rand.Int(rand.Reader, big.NewInt(MAX_ROLL+1))
	if err != nil {
		panic(err)
	}
	myRoll = nBig.Int64()
	conn.Privmsg(fmt.Sprintf("#%s", deployRequest.Swarm), fmt.Sprintf("[pc-roll] %d", myRoll))
}

func processChat(line, channel string) {
	msg := strings.Split(line, " ")
	switch msg[0] {
	case "[pc-deploy]":
		if len(msg[1]) == IPFS_LENGTH && len(msg[2]) == KEY_LENGTH {
			deployRequest = rpctest.Program{Hash: msg[1], Key: msg[2], Swarm: channel[1:]}
			rollForDeploy()
			// interpret as a pc deploy
			// log.Println("this was a correctly formatted pc deploy")
			// deployFromNetwork(rpctest.Program{Hash: msg[1], Key: msg[2]})
		}
	case "[pc-stop]":
		if len(msg[1]) == IPFS_LENGTH {
			log.Println("attempting to stop container")
			stopContainer(rpctest.Program{Hash: msg[1], Swarm: channel[1:]})
		}
	case "[pc-roll]":
		if roll, err := strconv.ParseInt(msg[1], 10, 64); err == nil && isDeploying {
			if roll >= MIN_ROLL && roll <= MAX_ROLL {
				log.Printf("received roll: %d\n", roll)
				if roll > highestRoll {
					// TODO: handle duplicate rolls / rolls with same value
					// (or just change to another form of leader election, but idk same same really unless one goes for
					// paxos or similar)
					highestRoll = roll
				}
			}
		}
	case "[pc-host]":
		if roll, err := strconv.ParseInt(msg[1], 10, 64); err == nil {
			log.Printf("highest roll was %d, my roll was %d", roll, myRoll)
			if roll == myRoll {
				log.Printf("i won the roll! deploying: %s\n", deployRequest.Hash)
				go deployFromNetwork(deployRequest)
			}
			myRoll = -1                       // reset roll
			deployRequest = rpctest.Program{} // empty deploy request
		}
	}
}

func RunDaemon(configpath string) {
	config = ReadConfig(configpath)
	readBookkeeping()
	// initialize RPC
	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+config.RPCPort)
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(Listener)
	rpc.Register(listener)
	go rpc.Accept(inbound)

	// intialize IRC
	conn = irc.IRC(config.Nick, config.Nick)
	// conn.VerboseCallbackHandler = true
	// conn.Debug = true
	conn.UseTLS = true
	conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	conn.AddCallback("001", func(event *irc.Event) {
		conn.Join(config.Channel)
		joinConnectedSwarms()
	})
	conn.AddCallback("366", func(event *irc.Event) {})
	conn.AddCallback("PRIVMSG", func(event *irc.Event) {
		log.Printf("%s: %s\t\t%s", event.Nick, event.Message(), event.Arguments[0])
		if strings.HasPrefix(event.Arguments[0], "#") {
			processChat(event.Message(), event.Arguments[0])
		}
		if event.Message() == "prompt" && strings.HasPrefix(event.Arguments[0], "#") {
			conn.Privmsg(event.Arguments[0], "boop")
		}
	})
	err = conn.Connect(config.Server)
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(1)
	}
	conn.Loop()
}
