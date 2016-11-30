package dockertest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

// static dockerfiles for the currently supported languages
const NODEJS = `FROM node:4-onbuild`
const PYTHON3 = `FROM python:3-onbuild
CMD [ "python", "./bot.py" ]`
const PYTHON2 = `FROM python:2-onbuild
CMD [ "python", "./bot.py" ]`

func listImages(client *docker.Client) {
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}

	for _, img := range imgs {
		fmt.Printf("id: %s\ncreated: %d\nsize: %d\nparent id: %s\n\n", img.ID, img.Created, img.Size, img.ParentID)
	}
}

func unpackTar(tarball, dest string) {
	cmd := exec.Command("tar", "xf", tarball, "-C", dest)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func CreateDockerfile(language, destination string) {
	dockerfile := ""
	switch language {
	case "py2":
		dockerfile = PYTHON2
	case "py3":
		dockerfile = PYTHON3
	case "nodejs":
		dockerfile = NODEJS
	default:
		panic(fmt.Sprintf("%s is currently an unsupported language!", language))
	}
	destination = strings.TrimRight(destination, "/")
	err := ioutil.WriteFile(fmt.Sprintf("%s/Dockerfile", destination), []byte(dockerfile), 0644)
	if err != nil {
		panic(err)
	}
}
func BuildImage(path, ipfshash string) {
	ipfshash = strings.ToLower(ipfshash)
	client := setupDocker()
	buildImage(client, path, ipfshash)
}

func buildImage(client *docker.Client, dockerfilepath string, ipfshash string) {
	outputbuf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:         ipfshash + "-image",
		OutputStream: outputbuf,
		ContextDir:   dockerfilepath,
	}
	if err := client.BuildImage(opts); err != nil {
		fmt.Printf("build image: %v\n", err)
		// log.Fatal(err)
		return
	}
	fmt.Println(outputbuf)
}

func RunContainer(ipfshash string) error {
	ipfshash = strings.ToLower(ipfshash)
	client := setupDocker()
	err := createContainer(ipfshash, client)
	if err != nil {
		fmt.Printf("create container: %v\n", err)
	}
	return startContainer(ipfshash, client)
}

func StopContainer(ipfshash string) {
	ipfshash = strings.ToLower(ipfshash)
	client := setupDocker()
	err := stopContainer(ipfshash, client)
	if err != nil {
		fmt.Printf("stop container: %v\n", err)
		return
	}
	removeContainer(ipfshash, client)
}

func createContainer(ipfshash string, client *docker.Client) error {
	dockerConfig := docker.Config{
		OpenStdin: true,
		Tty:       true,
		Image:     ipfshash + "-image",
	}
	_, err := client.CreateContainer(docker.CreateContainerOptions{
		Name:       ipfshash + "-container",
		Config:     &dockerConfig,
		HostConfig: nil,
	})
	if err != nil {
		return err
	}

	fmt.Println("successfully created the container")
	return nil
}

func startContainer(ipfshash string, client *docker.Client) error {
	err := client.StartContainer(ipfshash+"-container", nil)
	if err != nil {
		// panic(err)
		fmt.Printf("start container: %v\n", err)
		return err
	}
	fmt.Println("successfully started the container")
	return nil
}

func stopContainer(ipfshash string, client *docker.Client) error {
	err := client.StopContainer(ipfshash+"-container", 10)
	if err != nil {
		return err
	}
	fmt.Println("successfully stopped the container")
	return nil
}

func removeContainer(ipfshash string, client *docker.Client) {
	err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    ipfshash + "-container",
		Force: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully removed the container")
}

func setupDocker() *docker.Client {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	return client
}

// func main() {
// 	// endpoint := "unix:///var/run/docker.sock"
// 	// client, err := docker.NewClient(endpoint)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// listImages(client)
// 	// runContainer(client)
// 	// stopContainer(client)
// 	// createContainer(client)
// 	// path := "/home/cblgh/code/go/src/pc/dockertest/nice"
// 	// imagename := "untar-test2"
// 	// containername := "untar-container"
//
// 	// unpackTar("python-docker.tar.gz", path)
// 	// buildImage(client, path, imagename)
// 	// runContainer(imagename, containername, client)
// 	// stopContainer(containername, client)
// 	destination := "testaf"
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(dir)
// 	// destination = fmt.Sprintf("%s/%s", dir, destination)
// 	destination, err = filepath.Abs(destination)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(destination)
// 	// fmt.Println("creating dockerfile")
// 	// CreateDockerfile("nodejs", destination)
// 	// fmt.Println("creating tarball")
// 	// fmt.Printf("%s %s\n", destination, fmt.Sprintf("%s/pc-docker-setup.tar.gz", destination))
// 	// tartest.PackTar(destination, fmt.Sprintf("%s/pc-docker-setup.tar.gz", destination))
// 	// fmt.Println("finished!")
// }
