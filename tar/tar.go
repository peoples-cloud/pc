package tar

import (
	"os"
	"os/exec"
)

// TODO: find solution that doesn't depend on os environemnt i.e. entirely golang based solution for tar unpack & pack
func Unpack(tarball, dest string) {
	cmd := exec.Command("tar", "xf", tarball, "-C", dest, ".")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func Pack(src, tarball string) {
	fi, err := os.Stat(src)
	if err != nil {
		panic(err)
	}
	var cmd *exec.Cmd
	// determine whether packing a file or a directory
	switch mode := fi.Mode(); {
	case mode.IsDir():
		cmd = exec.Command("tar", "czf", tarball, "-C", src, ".")
	case mode.IsRegular():
		cmd = exec.Command("tar", "czf", tarball, src)
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

// func main() {
// 	Pack(os.Args[2], os.Args[1])
// }
