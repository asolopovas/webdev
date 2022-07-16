package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/asolopovas/dsync/webdev/lib"
	"github.com/txn2/txeh"
)

type Host struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type JsonConfig struct {
	Output   string `json:"output"`
	WorkDir  string `json:"workdir"`
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

	fileName := lib.PathResolve(conf.Output + host.Name + ".conf")
	templatePath := lib.PathResolve(conf.Template)
	template, err := ioutil.ReadFile(templatePath)

	nginxRoot := lib.AddTrailingSlash(conf.WorkDir) + host.Name
	if host.Type == "laravel" {
		nginxRoot = nginxRoot + "/public"
	}

	siteConfig := strings.Replace(string(template), "${APP_URL}", host.Name, -1)
	siteConfig = strings.Replace(siteConfig, "${NGINX_ROOT}", nginxRoot, -1)

	if host.Type == "wordpress" {
		siteConfig = strings.Replace(siteConfig, "# ${WORDPRESS}", "include                 extras/wordpress.conf;", -1)
	}

	err = ioutil.WriteFile(fileName, []byte(siteConfig), 0644)
	lib.ErrChk(err)

}

func HostsAdd(sites []Host) {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}
	hosts.AddHosts("127.0.0.1", []string{"phpmyadmin.test", "mariadb", "mailhog", "redis"})

	for _, site := range sites {
		fmt.Println("adding " + site.Name + " to /etc/hosts")
		hosts.AddHost("127.0.0.1", site.Name)
	}
	hfData := hosts.RenderHostsFile()

	// if you like to see what the outcome will
	// look like
	fmt.Println(hfData)

	hosts.Save()
}

func main() {
	var cFlag = flag.String("c", "webdev.json", "custom config path")
	var hFlag = flag.Bool("h", false, "Help")
	var hostsFlag = flag.Bool("hosts", false, "Help")
	flag.Parse()

	if *hFlag {
		flag.PrintDefaults()
		return
	}

	conf, err := getConfig(*cFlag)
	lib.ErrChk(err)

	if *hostsFlag {
		HostsAdd(conf.Hosts)
		return
	}

	lib.RmOldConfigs()
	for _, host := range conf.Hosts {
		generateHostConfig(host, conf)
	}

	lib.Cmd("docker", "compose build nginx app", true)
	lib.Cmd("docker", "compose up -d nginx app", false)

}
