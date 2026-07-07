package config

type ServerConfig struct {
	Host     string `mapstructure:"host"`
	HTTPPort int    `mapstructure:"http_port"`
	Port     int    `mapstructure:"port"`
}
