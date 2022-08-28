package config

import "os"

type Config struct {
	HCaptchaSecret string
}

func MakeConfig() Config {
	return Config{
		HCaptchaSecret: EnvOrDefault("HCAPTCHA_SECRET", "Fillthis"),
	}
}

func EnvOrDefault(envName string, defaultVal string) string {
	env := os.Getenv(envName)
	if env == "" {
		env = defaultVal
	}
	return env
}
