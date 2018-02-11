package config

import (
	"os"
	"strings"

	"github.com/grokify/gotilla/os/osutil"
	"github.com/joho/godotenv"
)

var DefaultPaths = []string{"ENV_PATH", ".env"}

func LoadDotEnv(paths ...string) error {
	if len(paths) == 0 {
		paths = DefaultPaths
	}

	envPathsSet := []string{}
	for _, envPathVar := range paths {
		envPath := strings.TrimSpace(os.Getenv(envPathVar))
		if len(envPath) > 0 {
			envPaths := strings.Split(envPath, ",")
			for _, envPath := range envPaths {
				envPath = strings.TrimSpace(envPath)
				if len(envPath) > 0 {
					exists, err := osutil.Exists(envPath)
					if err == nil && exists {
						envPathsSet = append(envPathsSet, envPath)
					}
				}
			}
		}
	}
	return godotenv.Load(envPathsSet...)
}
