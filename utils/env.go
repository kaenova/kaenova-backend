package utils

import "os"

func EnvOrDefault(envName string, defaultVal string) string {
	env := os.Getenv(envName)
	if env == "" {
		env = defaultVal
	}
	return env
}
