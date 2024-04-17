package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type DBConfig struct {
	Type      string `yaml:"connection"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Database  string `yaml:"database"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Collation string `yaml:"collation"`
	MinConn   int    `yaml:"min_connections"`
	MaxConn   int    `yaml:"max_connections"`
}

type JwtConfig struct {
	JWTSecret          string        `yaml:"jwt_secret"`
	JwtDuration        time.Duration `yaml:"jwt_duration"`
	JwtRefreshDuration time.Duration `yaml:"jwt_refresh_duration"`
}

type AppConfig struct {
	Port        string `yaml:"port"`
	Host        string `yaml:"host"`
	Env         string `yaml:"environment"`
	ProjectName string `yaml:"project_name"`
}

type Config struct {
	App AppConfig `yaml:"app_config"`

	JwtConfig JwtConfig `yaml:"jwt_config"`
	DBConfig  DBConfig  `yaml:"db_config"`
}

var config *Config

func Init(filepath string) error {

	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return nil
}

func GetConfig() *Config {
	return config
}
