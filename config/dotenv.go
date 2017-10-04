package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var EnvPathVars = []string{"ENV_PATH", "."}

func LoadDotEnv() error {
	envPathsSet := []string{}
	for _, envPathVar := range EnvPathVars {
		envPath := strings.TrimSpace(os.Getenv(envPathVar))
		if len(envPath) > 0 {
			envPaths := strings.Split(envPath, ",")
			for _, envPath := range envPaths {
				envPath = strings.TrimSpace(envPath)
				if len(envPath) > 0 {
					envPathsSet = append(envPathsSet, envPath)
				}
			}
		}
	}
	return godotenv.Load(envPathsSet...)
}
