package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"
	"os"

	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
)

type BridgeConfig struct {
	Service Service
}

type Service struct {
	Name    string
	Port    int
	Address string
}

func parseConfig(config_file string) BridgeConfig {
	if _, err := os.Stat(config_file); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "configuration file  -%s do not exists\n", config_file)
		os.Exit(2)
	}

	file, e := ioutil.ReadFile(config_file)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(2)
	}

	var config BridgeConfig
	err := yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}

func main() {

	var config_file = flag.String("config_file", "", "Configuration file for the bridge")
	required := []string{"config_file"}
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required -%s argument\n", req)
			os.Exit(2)
		}
	}

	config := parseConfig(*config_file)
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		fmt.Println(err)
	}

	agent := client.Agent()
	service := &api.AgentServiceRegistration{ID: config.Service.Name,
		Name:    config.Service.Name,
		Tags:    nil,
		Port:    config.Service.Port,
		Address: config.Service.Address}

	for {
		agent.ServiceRegister(service)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(20*time.Second)

	}
}
