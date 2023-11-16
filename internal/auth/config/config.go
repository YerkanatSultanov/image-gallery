package config

import "time"

type Config struct {
	Auth      Auth
	Server    Server    `yaml:"Server"`
	Database  Database  `yaml:"Database"`
	Transport Transport `yaml:"Transport"`
}

type Auth struct {
	secretKey string `yaml:"secretKey"`
}

type UserGrpcTransport struct {
	Host string `yaml:"host"`
}

type Server struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  string `yaml:"sslMode"`
}

type Transport struct {
	User     UserTransport     `yaml:"user"`
	UserGrpc UserGrpcTransport `yaml:"userGrpcTransport"`
}

type UserTransport struct {
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
}
