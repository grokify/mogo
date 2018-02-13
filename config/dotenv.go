package config

import (
	"os"
	"strings"

	"github.com/grokify/gotilla/os/osutil"
	"github.com/joho/godotenv"
)

const (
	EnvPathVar = "ENV_PATH"
	LocalPath  = "./.env"
)

var DefaultPaths = []string{os.Getenv(EnvPathVar), "./.env"}

func LoadEnvDefaults() error {
	envPathsSet := []string{}

	for _, defaultPath := range DefaultPaths {
		exists, err := osutil.Exists(defaultPath)
		if err == nil && exists {
			envPathsSet = append(envPathsSet, defaultPath)
		}
	}

	if len(envPathsSet) > 0 {
		return godotenv.Load(envPathsSet...)
	}
	return godotenv.Load()
}

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
