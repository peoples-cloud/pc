package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"unicode"

	"github.com/peoples-cloud/pc/daemon"
	"github.com/peoples-cloud/pc/rpc"

	"github.com/spf13/cobra"
)

var flagvar int
var swarmName, swarmPassword, configpath, rpcport string
var requirePass bool
var noPass = false
var noHost = false

func sendRPC(command string, info rpcmsg.Info) rpcmsg.Message {
	client := dialDaemon()
	var reply rpcmsg.Message
	command = Capitalize(command)
	err := client.Call(fmt.Sprintf("Listener.%s", command), info, &reply)
	if err != nil {
		log.Fatal(err)
	}
	return reply
}

func Capitalize(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func checkArgLength(n int, args []string, cmd *cobra.Command) {
	switch {
	case len(args) == 0:
		cmd.Help()
		os.Exit(0)
	case len(args) < n:
		fmt.Println(cmd.Name() + ": too few arguments")
		fmt.Println(cmd.Use)
		os.Exit(0)
	case len(args) > n:
		fmt.Println(cmd.Name() + ": too many arguments")
		fmt.Println(cmd.Use)
		os.Exit(0)
	}
}

func main() {
	// TODO: find solution to use rpcport from config, or just remove rpcport from config

	// NOTE: remove cmdTest, maybe?
	var cmdDaemon = &cobra.Command{
		Use:     "daemon <path to config>",
		Short:   "Starts the pc daemon",
		Long:    `Starts the pc daemon that connects to previously saved swarms, and allows you to host & deploy programs`,
		Example: `pc daemon config.toml`,

		Run: func(cmd *cobra.Command, args []string) {
			checkArgLength(1, args, cmd)
			configpath, err := filepath.Abs(args[0])
			if err != nil {
				fmt.Printf("pc daemon: %v\n", err)
				os.Exit(1)
			}
			daemon.RunDaemon(configpath)
		},
	}

	// NOTE: remove cmdTest, maybe?
	// var cmdTest = &cobra.Command{
	// 	Use:     "test <hash> <key>",
	// 	Short:   "Test-deploy a program on this node",
	// 	Long:    `Test-deploy a program on this node`,
	// 	Example: `pc test QmWLwhxjJuqff7cXniyZegg2NQnYFhVp5XcyyeLRjRqBtu vIhCg78Ef0hxGfvpEIeONTpDJGIL2UQWPh0fH8nTgzs`,
	//
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		checkArgLength(2, args, cmd)
	// 		fmt.Printf("doing testdeploy with\nhash: %s\nkey:%s", args[0], args[1])
	// 		info := rpctest.RPCInfo{Hash: args[0], Key: args[1]}
	// 		_ = sendRPC(cmd.Name(), info)
	// 	},
	// }
	/* cobra cli command stuff */
	var cmdDeploy = &cobra.Command{
		Use:   "deploy <swarm> <path> <language>",
		Short: "Deploy a program to be run by a node in the swarm",
		Long: `Deploys a program to the specified swarm using a combination of technologies
            `,
		Example: `pc deploy peoples-swarm ~/bots/cool-bot py2

currently available language options are:
* py2
* py3
* nodejs`,
		Run: func(cmd *cobra.Command, args []string) {
			checkArgLength(3, args, cmd)
			swarm := args[0]
			path := args[1]
			lang := args[2]

			path, err := filepath.Abs(path)
			if err != nil {
				panic(err)
			}
			fmt.Printf("deploying %s to %s...\n", path, swarm)
			info := rpcmsg.Info{Swarm: swarm, Path: path, Language: lang}
			reply := sendRPC(cmd.Name(), info)
			fmt.Printf("deploy:\n%s\n", reply.Msg)
		},
	}

	var cmdStop = &cobra.Command{
		Use:     "stop <swarm> <hash>",
		Short:   "Stop a deployed program",
		Long:    `Halts execution of a previously deployed program`,
		Example: "pc stop peoples-swarm <hash>",
		Run: func(cmd *cobra.Command, args []string) {
			checkArgLength(2, args, cmd)
			swarm := args[0]
			hash := args[1]
			info := rpcmsg.Info{Swarm: swarm, Hash: hash}
			_ = sendRPC(cmd.Name(), info)
			fmt.Printf("%s: stopped %s\n", swarm, hash)
		},
	}

	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create a new swarm",
		Long: `Creates a new peoples cloud swarm. 
If no options are provided the generated swarm name is returned to stdout`,
		Example: `pc create -n peoples-swarm`,
		// pc create -n peoples-swarm -p forthepeople
		// pc create -n peoples-swarm --no-pw`,
		Run: func(cmd *cobra.Command, args []string) {
			// checkArgLength(0, args, cmd)
			if len(args) > 0 {
				fmt.Println("error: create takes swarm name as a flag:\npc create -n <swarm name>")
				os.Exit(0)
			}
			info := rpcmsg.Info{Swarm: swarmName, Password: swarmPassword}
			reply := sendRPC(cmd.Name(), info)
			fmt.Printf("created %s\n", reply.Msg)
		},
	}

	var cmdJoin = &cobra.Command{
		Use:   "join <swarm>",
		Short: "Join a new swarm",
		Long: `Joins a new pc swarm 
		`,
		Run: func(cmd *cobra.Command, args []string) {
			checkArgLength(1, args, cmd)
			swarmName = args[0]
			info := rpcmsg.Info{Swarm: swarmName, Password: swarmPassword}
			_ = sendRPC(cmd.Name(), info)
			fmt.Printf("joined %s\n", swarmName)
		},
	}

	var cmdLeave = &cobra.Command{
		Use:   "leave <swarm>",
		Short: "Leave a swarm",
		Long: `Leaves a previously joined swarm 
		`,
		Run: func(cmd *cobra.Command, args []string) {
			checkArgLength(1, args, cmd)
			swarmName = args[0]
			info := rpcmsg.Info{Swarm: swarmName}
			_ = sendRPC(cmd.Name(), info)
			fmt.Printf("left %s\n", swarmName)
		},
	}

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "Lists all connected swarms and deployed programs",
		Long: `Lists all currently connected swarms, your deployed programs and their corresponding hashes and keys
		`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("swarm\t\tprogram\t\thash\t\tkey\n")
			fmt.Printf("%s\t\t%s\t\t%s\t\t%s\n", "peoples-cloud", "~/mybots/cyberbot", "QmQyy6ZvLnqXUDyRMcu5xvGD6guhbZ7ydKQTPt4NqkMv9W", "DGZPpulZVFLolSFc9BLbMBYfqse4_NgJSxl7PSbM_Bk=")
		},
	}

	cmdCreate.Flags().StringVarP(&swarmName, "name", "n", "", "swarm name")
	// not implemented yet:
	// 	passwords for channels
	// 	no-hosting option
	// cmdCreate.Flags().StringVarP(&swarmPassword, "password", "p", "", "supply the password required to join swarm")
	// cmdCreate.Flags().BoolVar(&noPass, "nopw", false, "disable password for swarm")

	// cmdJoin.Flags().StringVarP(&swarmPassword, "password", "p", "", "supply the password required to join swarm")
	// cmdJoin.Flags().BoolVar(&noHost, "no-host", false, "disable hosting for this node")

	var rootCmd = &cobra.Command{Use: "pc"}
	rootCmd.AddCommand(cmdDeploy, cmdCreate, cmdJoin, cmdLeave, cmdList, cmdStop /* cmdTest,*/, cmdDaemon)
	rootCmd.PersistentFlags().StringVarP(&rpcport, "rpcport", "r", "42586", "specifies the port used to communicate with the pc daemon using rpc, should match rpcport in the config")
	rootCmd.Execute()

}

func dialDaemon() *rpc.Client {
	client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%s", rpcport))
	if err != nil {
		fmt.Println("run pc daemon in another window before trying to deploy")
		log.Fatal(err)
	}
	return client
}
