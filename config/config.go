package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiUrl      string
	ApiOrg      string
	ApiSecret   string
	Port        string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		ApiUrl:        cfg.Section("gpt").Key("api_url").String(),
		ApiOrg:        cfg.Section("gpt").Key("api_org").String(),
		ApiSecret:     cfg.Section("gpt").Key("api_secret").String(),
		Port:          cfg.Section("web").Key("port").String(),
	}
}
