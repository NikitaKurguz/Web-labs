package store

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DataBase ConfigDB `toml:"database"`
}

type ConfigDB struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DB_Name  string `toml:"dbname"`
	SSL_Mode string `toml:"sslmode"`
}

func NewConfig() *Config {
	return &Config{
		DataBase: ConfigDB{
			Host:     "localhost",
			Port:     "5440",
			User:     "postgres",
			Password: "postgres",
			DB_Name:  "primarte",
			SSL_Mode: "disable",
		},
	}
}

func (c *Config) Load(path string) error {
	_, err := toml.DecodeFile(path, c)

	if v := os.Getenv("DB_HOST"); v != "" {
		c.DataBase.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		c.DataBase.Port = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		c.DataBase.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		c.DataBase.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		c.DataBase.DB_Name = v
	}
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		c.DataBase.SSL_Mode = v
	}
	return err
}

func (c *ConfigDB) ConnectionString() string {
	return "host=" + c.Host + " port=" + c.Port +
		" user=" + c.User + " password=" + c.Password +
		" dbname=" + c.DB_Name + " sslmode=" + c.SSL_Mode
}
