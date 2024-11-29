package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth AuthConfig `yaml:"auth"`
	PasswordHashing PasswordHashingConfig `yaml:"passwordHashing"`
}

func LoadConfig(filepath string) (*Config, error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	DriverName string `yaml:"driverName"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Name string `yaml:"name"`
	SslMode string `yaml:"sslmode"`
	TimeZone string `yaml:"timeZone"`
}

type AuthConfig struct {
	RefreshToken JwtTokenConfig `yaml:"refreshToken"`
	AccessToken JwtTokenConfig `yaml:"accessToken"`
	Issuer string `yaml:"issuer"`
}

type JwtTokenConfig struct {
	Duration time.Duration `yaml:"duration"`
	SigningKey string `yaml:"signingKey"`
}

type PasswordHashingConfig struct {
	Algorithm string `yaml:"algorithm"`
	Bcrypt BcryptConfig `yaml:"bcrypt"`
}

type BcryptConfig struct {
	Cost int `yaml:"cost"`
}