package apiserver

import "github.com/BurntSushi/toml"

type Config struct {
	Server ServerConfig `toml:"server"`
}

type ServerConfig struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			BindAddr: ":8080",
			LogLevel: "debug",
		},
	}
}

func (c *Config) Load(path string) error {
	_, err := toml.DecodeFile(path, c)
	return err
}
