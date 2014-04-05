package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type dbConfiguration struct {
	Name string
	User string
	Pass string
	Pool int
}

type configuration struct {
	UploadPath string
	Db         dbConfiguration
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

	return config
}
