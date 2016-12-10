package ipfstest

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func IPFSAdd(filepath string) string {
	out, err := exec.Command("ipfs", "add", filepath).Output()
	if err != nil {
		log.Fatalf("ipfs add: %v\n", err)
	}
	hash := strings.Split(string(out), " ")[1]
	return hash
}

func IPFSGet(hash, dest string) {
	fmt.Printf("hash: %s, dest: %s\n", hash, dest)
	cmd := exec.Command("ipfs", "get", hash, "-o", dest)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("ipfs get: %v\n", err)
	}
}

// func main() {
// 	// ipfsGet(os.Args[1], os.Args[2])
// 	hash := ipfsAdd(os.Args[1])
// 	fmt.Printf("hash: %s\n", hash)
// }
