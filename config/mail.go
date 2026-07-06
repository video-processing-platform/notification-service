package config

type MailConfig struct {
	Provider string

	From string

	MailHog  MailHogConfig
	MailTrap MailTrapConfig
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
