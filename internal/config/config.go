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
	filepath string
	Jira     struct {
		BaseURL   string `yaml:"base_url"`
		UserEmail string `yaml:"user_email"`
		APIToken  string `yaml:"api_token"`
		Project   string `yaml:"project"`
	} `yaml:"jira"`
}

func Load() (*Config, error) {
	configDir, err := getLocalConfigDir()
	if err != nil {
		return nil, err
	}
	configFilepath := filepath.Join(configDir, "config.yml")
	config, err := readConfig(configFilepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, io.EOF) {
			config, err = newConfig(configFilepath)
			if err != nil {
				return nil, err
			}
			if err = config.Save(); err != nil {
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

func readConfig(filepath string) (*Config, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	conf := &Config{
		filepath: filepath,
	}
	yamlDecoder := yaml.NewDecoder(file)
	if err = yamlDecoder.Decode(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func newConfig(filepath string) (*Config, error) {
	config := &Config{
		filepath: filepath,
	}
	if err := config.Init(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Save() error {
	file, err := os.OpenFile(c.filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	yamlEncoder := yaml.NewEncoder(file)
	defer yamlEncoder.Close()
	yamlEncoder.SetIndent(2)
	return yamlEncoder.Encode(c)
}
