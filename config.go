package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var configfile string = "/etc/conduit.conf"

type Config struct {
	Backends    string
	Bind        string
	MonitorBind string
	Mode        string
	Certfile    string
	Keyfile     string
}

//read config file
func ReadConfig() Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file not found: ", configfile)
	}
	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
