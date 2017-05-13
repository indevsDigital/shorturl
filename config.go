package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Httpport   string `json:"httpport"`
	DBUser     string `json:"db_user"`
	DBName     string `json:"db_name"`
	DBHost     string `json:"db_host"`
	DBPassword string `json:"db_password"`
}

func NewConfig() *Config {
	var c Config
	err := c.LoadConfig()
	if err != nil {
		panic(err)
	}
	return &c
}

func (conf *Config) LoadConfig() error {
	w, err := ioutil.ReadFile("conf.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(w, conf)
	if err != nil {
		return err
	}
	return nil
}
