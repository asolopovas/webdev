package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RmOldConfigs() {
	dirPath := os.Getenv("HOME") + "/www/dev/nginx/sites"
	files, err := ioutil.ReadDir(dirPath)
	ErrChk(err)

	for _, file := range files {
		if file.Name() == ".gitkeep" {
			continue
		}
		filePath := filepath.Join(dirPath, file.Name())
		err := os.Remove(filePath)
		ErrChk(err)
	}

}

func Cmd(command string, arguments string, silent bool) {
	os.Chdir(os.Getenv("HOME") + "/www/dev")
	args := strings.Split(arguments, " ")
	fmt.Println(args)
	fmt.Println(command)
	cmd := exec.Command(command, args...)
	if !silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	cmd.Run()
}
