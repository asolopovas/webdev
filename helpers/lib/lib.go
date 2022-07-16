package lib

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ErrChk(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PathResolve(path string) string {
	dirname, _ := os.UserHomeDir()
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(dirname, path[2:])
	}
	return filepath.Join(dirname, path)
}
