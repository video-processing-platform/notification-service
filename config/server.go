package config

type ServerConfig struct {
	Host     string `mapstructure:"host"`
	HttpPort int    `mapstructure:"http_port"`
	Port     int    `mapstructure:"port"`
}
