package config

type Config struct {
	HCaptchaSecret string
}

func MakeConfig(c Config) Config {
	return c
}
