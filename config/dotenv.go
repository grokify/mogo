package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	iom "github.com/grokify/simplego/io/ioutilmore"
	"github.com/grokify/simplego/os/osutil"
	"github.com/grokify/simplego/type/stringsutil"
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

func LoadDotEnv(paths ...string) ([]string, error) {
	return LoadDotEnvSkipEmptyInfo(paths...)
}

func LoadDotEnvSkipEmptyInfo(paths ...string) ([]string, error) {
	if len(paths) == 0 {
		paths = DefaultPaths()
	}

	envPaths := iom.FilterFilenamesSizeGtZero(paths...)
	envPaths = stringsutil.Dedupe(envPaths)

	if len(envPaths) > 0 {
		return envPaths, godotenv.Load(envPaths...)
	}
	return envPaths, nil
}

func LoadDotEnvSkipEmpty(paths ...string) error {
	_, err := LoadDotEnvSkipEmptyInfo(paths...)
	return err
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
		isFile, err := osutil.IsFileWithSizeGtZero(fixedPath)
		if err != nil {
			return fixedPath, err
		} else if !isFile {
			return fixedPath, fmt.Errorf("Path is not a file or is 0 size [%v]", fixedPath)
		}
		return fixedPath, nil
	}

	envPath = strings.TrimSpace(envPath)
	if len(envPath) > 0 {
		isFile, err := osutil.IsFileWithSizeGtZero(fixedPath)
		if err != nil {
			return envPath, err
		} else if !isFile {
			return envPath, fmt.Errorf("Path is not a file or is 0 size [%v]", envPath)
		}
		return envPath, nil
	}

	thisDirPath := "./.env"
	isFile, err := osutil.IsFileWithSizeGtZero(thisDirPath)
	if err != nil {
		return thisDirPath, err
	} else if !isFile {
		return thisDirPath, fmt.Errorf("Path is not a file or is 0 size [%v]", thisDirPath)
	}
	return thisDirPath, nil
}
