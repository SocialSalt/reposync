package main

import (
	"os"
	"path/filepath"

	xdg "github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

type Repo struct {
	Name     string   `yaml:"name,omitempty"`
	Excludes []string `yaml:"excludes,omitempty"`
}

type Config struct {
	Debug          bool     `yaml:"debug"`
	GlobalExcludes []string `yaml:"global_excludes,omitempty"`
	Repos          []Repo   `yaml:"repos,omitempty"`
}

func defaultExcludes() []string {
	return []string{".git", "target", "__pycache__", ".DS_Store"}
}

func loadConfig() (Config, error) {
	configFile := filepath.Join(xdg.ConfigHome, "reposync/reposync.yaml")

	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		cfg := Config{
			GlobalExcludes: defaultExcludes(),
		}
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return cfg, err
		}
		err = os.Mkdir(filepath.Dir(configFile), 0755)
		if err != nil {
			return cfg, err
		}
		err = os.WriteFile(configFile, data, 0755)
		if err != nil {
			return cfg, err
		}
		return cfg, nil

	} else {
		yamlFile, err := os.ReadFile(configFile)
		if err != nil {
			return Config{}, err
		}
		cfg := Config{}
		err = yaml.Unmarshal(yamlFile, &cfg)
		if err != nil {
			return Config{}, err
		}
		return cfg, nil
	}
}
