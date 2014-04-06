package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

type dbConfiguration struct {
	Name string
	User string
	Pass string
	Port int
	Pool int
}

type configuration struct {
	UploadPath    string
	ServerAddress string
	User          string
	Password      string
	Db            dbConfiguration
}

func (c *configuration) absolutePath(path string) string {
	return filepath.Join(c.UploadPath, path)
}

func loadConfiguration(configPath string) configuration {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config configuration
	if err = json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	if config.User == "" || config.Password == "" {
		log.Fatal("user and password are required in the config")
	}

	return config
}
