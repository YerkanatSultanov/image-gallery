package config

type Config struct {
	Server     Server     `yaml:"Server"`
	Database   Database   `yaml:"Database"`
	GrpcServer GrpcServer `yaml:"GrpcServer"`
	Kafka      Kafka      `yaml:"Kafka"`
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

type GrpcServer struct {
	Port string `yaml:"Port"`
}

type Kafka struct {
	Brokers  []string `yaml:"brokers"`
	Producer Producer `yaml:"producer"`
	Consumer Consumer `yaml:"consumer"`
}

type Producer struct {
	Topic string `yaml:"topic"`
}

type Consumer struct {
	Topics []string `yaml:"topics"`
}
