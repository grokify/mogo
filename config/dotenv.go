package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	iom "github.com/grokify/gotilla/io/ioutilmore"
	"github.com/grokify/gotilla/os/osutil"
	"github.com/joho/godotenv"
)

var (
	EnvPathVar = "ENV_PATH"
	LocalPath  = "./.env"
)

func DefaultPaths() []string {
	return []string{os.Getenv(EnvPathVar), LocalPath}
}

func LoadEnvDefaults() error {
	envPathsSet := []string{}

	for _, defaultPath := range DefaultPaths() {
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

func LoadDotEnvSkipEmpty(paths ...string) error {
	if len(paths) == 0 {
		paths = DefaultPaths()
	}

	envPaths := iom.FilterFilenamesSizeGtZero(paths...)

	if len(envPaths) > 0 {
		return godotenv.Load(envPaths...)
	}
	return nil
}

func LoadDotEnvFirst(paths ...string) error {
	if len(paths) == 0 {
		paths = DefaultPaths()
	}

	envPaths := iom.FilterFilenamesSizeGtZero(paths...)

	if len(envPaths) > 0 {
		return godotenv.Load(envPaths[0])
	}
	return nil
}

// GetDotEnvVal retrieves a single var from a `.env` file path
func GetDotEnvVal(envPath, varName string) (string, error) {
	cmd := fmt.Sprintf("grep %s '%s' | rev | cut -d= -f1 | rev", varName, envPath)

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("Failed to execute command: %s", cmd)
	}
	return string(out), nil
}

// LoadEnvPaths attempts to load an explicit, env and current path.
// It returns an error ifexplicit/env paths are present but do not
// exist or are empty. This is was designed to flexibly handle common
// .env file paths in a prioritzed and differentiated order.
func LoadEnvPathsPrioritized(fixedPath, envPath string) error {
	if goodPath, err := checkEnvPathsPrioritized(fixedPath, envPath); err != nil {
		return err
	} else if len(goodPath) > 0 {
		return godotenv.Load(goodPath)
	}
	return nil
}

func checkEnvPathsPrioritized(fixedPath, envPath string) (string, error) {
	fixedPath = strings.TrimSpace(fixedPath)
	if len(fixedPath) > 0 {
		return fixedPath, iom.IsFileWithSizeGtZero(fixedPath)
	}

	envPath = strings.TrimSpace(envPath)
	if len(fixedPath) > 0 {
		return envPath, iom.IsFileWithSizeGtZero(envPath)
	}

	thisDirPath := "./.env"
	err := iom.IsFileWithSizeGtZero(thisDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return thisDirPath, err
		}
	}
	return thisDirPath, nil
}
