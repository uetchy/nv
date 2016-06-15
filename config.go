package main

import (
	"encoding/json"
	"github.com/Songmu/prompter"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
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

	err = ioutil.WriteFile(destPath, blob, 0600)
	if err != nil {
		return err
	}

	return nil
}

func generateConfig() (err error) {
	email := prompter.Prompt("Email", "")
	viper.Set("Email", email)
	password := prompter.Prompt("Password", "")
	viper.Set("Password", password)

	defaultConfigPath, _ := homedir.Expand("~/.config/nv/config.json")
	err = os.MkdirAll(filepath.Dir(defaultConfigPath), 0755)
	if err != nil {
		return err
	}

	err = saveConfig(defaultConfigPath)
	if err != nil {
		return err
	}
	return nil
}
