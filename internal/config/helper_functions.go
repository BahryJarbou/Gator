package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	completePath := filepath.Join(homePath + fileName)

	return completePath, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	configfile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer configfile.Close()

	encoder := json.NewEncoder(configfile)

	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
