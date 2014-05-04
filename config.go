package main

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
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
	Db            dbConfiguration
	SessionSecret string
	ClientId      string // auth.sch client id
	ClientSecret  string // auth.sch client secret
	Secure        bool   // using https
	GroupId       int    // id of the group that is allowed to upload images. if <= 0 everybody allowed
}

func (c *configuration) absolutePath(path string) string {
	return filepath.Join(c.UploadPath, path)
}

func (c *configuration) oauth() *oauth.Config {
	var redirect bytes.Buffer

	redirect.WriteString("http")
	if c.Secure {
		redirect.WriteString("s")
	}
	redirect.WriteString("://")
	redirect.WriteString(strings.TrimRight(c.ServerAddress, "/"))
	redirect.WriteString("/auth")

	return &oauth.Config{
		ClientId:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Scope:        "basic eduPersonEntitlement",
		TokenURL:     "https://auth.sch.bme.hu/oauth2/token",
		AuthURL:      "https://auth.sch.bme.hu/site/login",
		RedirectURL:  redirect.String(),
	}
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
