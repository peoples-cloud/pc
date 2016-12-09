package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func MakeDir(path, hash string) string {
	path, err := filepath.Abs(path)
	folder := fmt.Sprintf("%s/%s", path, hash)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s", path, hash), os.FileMode(0777))
	if err != nil {
		panic(err)
	}
	return folder
}
