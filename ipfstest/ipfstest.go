package ipfstest

import (
	"fmt"
	"os/exec"
	"strings"
)

func IPFSAdd(filepath string) string {
	out, err := exec.Command("ipfs", "add", filepath).Output()
	if err != nil {
		panic(err)
	}
	hash := strings.Split(string(out), " ")[1]
	return hash
}

func IPFSGet(hash, dest string) {
	cmd := exec.Command("ipfs", "get", hash)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ipfs get: %v", err)
		panic(err)
	}
	cmd = exec.Command("mv", hash, dest)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

// func main() {
// 	// ipfsGet(os.Args[1], os.Args[2])
// 	hash := ipfsAdd(os.Args[1])
// 	fmt.Printf("hash: %s\n", hash)
// }
