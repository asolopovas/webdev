package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/asolopovas/webdev/helpers/lib"
)

type Host struct {
	Host     string `json:"host"`
	Template string `json:"template"`
	WorkDir  string `json:"workdir"`
	Type     string `json:"type"`
}

type JsonConfig struct {
	Output   string `json:"output"`
	Template string `json:"template"`
	Hosts    []Host `json:"hosts"`
}

func getConfig(configPath string) (JsonConfig, error) {
	result := JsonConfig{}

	config, err := ioutil.ReadFile(configPath)

	json.Unmarshal([]byte(config), &result)
	return result, err
}

func generateHostConfig(host Host, conf JsonConfig) {

	fileName := lib.PathResolve(conf.Output + host.Host + ".conf")
	templatePath := lib.PathResolve(conf.Template)

	template, err := ioutil.ReadFile(templatePath)

	siteConfig := strings.Replace(string(template), "${APP_URL}", host.Host, -1)
	if host.Type == "wordpress" {
		siteConfig = strings.Replace(siteConfig, "${NGINX_ROOT}", host.WorkDir, -1)
		siteConfig = strings.Replace(siteConfig, "# ${WORDPRESS}", "include                 extras/wordpress.conf;", -1)
	} else {
		siteConfig = strings.Replace(siteConfig, "${NGINX_ROOT}", host.WorkDir+"/public", -1)
	}

	err = ioutil.WriteFile(fileName, []byte(siteConfig), 0644)
	lib.ErrChk(err)

}

func main() {
	var cFlag = flag.String("c", "webdev.json", "custom config path")
	flag.Parse()

	conf, err := getConfig(*cFlag)
	if err != nil {
		fmt.Println(err)
	}

	for _, host := range conf.Hosts {
		generateHostConfig(host, conf)
	}

}
