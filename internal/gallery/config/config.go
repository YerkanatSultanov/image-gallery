package config

type Config struct {
	Server    Server    `yaml:"Server"`
	Database  Database  `yaml:"Database"`
	Transport Transport `yaml:"Transport"`
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

type AuthGrpcTransport struct {
	Host string `yaml:"host"`
}

type UserGrpcTransport struct {
	Host string `yaml:"host"`
}

type Transport struct {
	AuthGrpc AuthGrpcTransport `yaml:"galleryGrpc"`
	UserGrpc UserGrpcTransport `yaml:"userGrpc"`
}
