package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Jira struct {
		BaseURL   string `yaml:"base_url"`
		UserEmail string `yaml:"user_email"`
		APIToken  string `yaml:"api_token"`
	} `yaml:"jira"`
}

func Load() (*Config, error) {
	configDir, err := getLocalConfigDir()
	if err != nil {
		return nil, err
	}
	configFile, err := getConfigFile(configDir)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	config, err := readConfig(configFile)
	if err != nil {
		if errors.Is(err, io.EOF) {
			config = &Config{}
			if err = initConfig(config); err != nil {
				return nil, err
			}
			if err = saveConfig(configFile, config); err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, err
	}

	return config, nil
}

func getLocalConfigDir() (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(workdir, ".jgitra")
	info, err := os.Stat(configDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err = os.Mkdir(configDir, 0755); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else if !info.IsDir() {
		return "", fmt.Errorf("'%s' is not a directory", configDir)
	}
	return configDir, nil
}

func getConfigFile(configDir string) (*os.File, error) {
	configFile := filepath.Join(configDir, "config.yml")
	file, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func readConfig(configFile *os.File) (*Config, error) {
	conf := &Config{}
	yamlDecoder := yaml.NewDecoder(configFile)
	if err := yamlDecoder.Decode(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func saveConfig(configFile *os.File, config *Config) error {
	yamlEncoder := yaml.NewEncoder(configFile)
	defer yamlEncoder.Close()
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(config); err != nil {
		return err
	}
	return nil
}
