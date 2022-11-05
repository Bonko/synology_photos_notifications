package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Rootpath string `yaml:"rootpath"`
	Users    []User `yaml:"users"`
}
type User struct {
	name  string `yaml:"name"`
	email string `yaml:"email"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decodes
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
