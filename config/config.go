package config

import (
    "gopkg.in/yaml.v3"
    "os"
)

type Config struct {
    Pinger App `yaml:"pinger"`
}

type App struct {
    Period int    `yaml:"period"`
    Api    string `yaml:"api"`
}

func ParseConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    var config Config
    err = yaml.NewDecoder(file).Decode(&config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}
