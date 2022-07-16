package lib

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func AddTrailingSlash(str string) string {
	match, _ := regexp.MatchString("/$", str)
	if match {
		return str
	}
	return str + "/"
}
