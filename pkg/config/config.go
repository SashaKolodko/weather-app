package config

import (
    "io"
    "gopkg.in/yaml.v3"
)

type Provider struct {
    Type string `yaml:"type"`
}

type Location struct {
    Lat  float64 `yaml:"lat"`
    Long float64 `yaml:"long"`
}

type Config struct {
    P Provider `yaml:"provider"`
    L Location `yaml:"location"`
}

type ConfigFile struct {
    C Config `yaml:"service"`
}

func Parse(r io.Reader) (Config, error) {
    var c ConfigFile
    if err := yaml.NewDecoder(r).Decode(&c); err != nil {
        return Config{}, err
    }
    return c.C, nil
}