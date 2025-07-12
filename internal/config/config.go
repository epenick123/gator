package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	home_directory, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	config_file_path := filepath.Join(home_directory, ".gatorconfig.json")

	json_file, err := os.Open(config_file_path)
	if err != nil {
		return Config{}, err
	}
	defer json_file.Close()

	byte_value, err := io.ReadAll(json_file)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(byte_value, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username

	home_directory, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	config_file_path := filepath.Join(home_directory, ".gatorconfig.json")

	json_file, err := os.Create(config_file_path)
	if err != nil {
		return err
	}
	defer json_file.Close()

	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	_, err = json_file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
