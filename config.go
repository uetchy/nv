package main

import (
	"encoding/json"
	"github.com/Songmu/prompter"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io/ioutil"
)

type config struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/nv")
	err := viper.ReadInConfig()
	return err
}

func saveConfig(destPath string) error {
	var C config
	err := viper.Unmarshal(&C)
	if err != nil {
		return err
	}

	blob, err := json.Marshal(C)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destPath, blob, 0644)
	if err != nil {
		return err
	}

	return nil
}

func generateConfig() error {
	email := prompter.Prompt("Email", "")
	viper.Set("Email", email)
	password := prompter.Prompt("Password", "")
	viper.Set("Password", password)
	defaultConfigPath, _ := homedir.Expand("~/.config/nv/config.json")
	saveConfig(defaultConfigPath)
	return nil
}
