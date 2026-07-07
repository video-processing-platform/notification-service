package config

type MailConfig struct {
	Provider      string
	From          string
	Timeout       int `mapstructure:"timeout"`
	RetryAttempts int `mapstructure:"retry_attempts"`
	RetryDelay    int `mapstructure:"retry_delay"`
	MailHog       MailHogConfig
	MailTrap      MailTrapConfig
}

type MailHogConfig struct {
	Host string
	Port int
}

type MailTrapConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}
