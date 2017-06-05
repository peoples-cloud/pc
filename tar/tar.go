package tar

import "os/exec"

// TODO: find solution that doesn't depend on os environemnt i.e. entirely golang based solution for tar unpack & pack
func Unpack(tarball, dest string) {
	cmd := exec.Command("tar", "xf", tarball, "-C", dest, ".")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func Pack(src, tarball string) {
	cmd := exec.Command("tar", "czf", tarball, "-C", src, ".")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// func main() {
// 	packTar(os.Args[1], os.Args[2])
// }
