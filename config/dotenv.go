package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/type/slicesutil"
	flags "github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

var (
	EnvPathVar = "ENV_PATH"
	LocalPath  = "./.env"
)

func DefaultPaths() []string {
	return []string{os.Getenv(EnvPathVar), LocalPath}
}

// LoadDotEnv loads a set of `.env` files given the supplied path.
// If no files are supplied, it will load the files from `DefaultPaths()`.
// A maximum of `n` files will be loaded. If `n` is 0 or less, all files
// will be loaded.
func LoadDotEnv(paths []string, n int) ([]string, error) {
	if len(paths) == 0 {
		paths = DefaultPaths()
	}

	envPaths := osutil.FilenamesFilterSizeGTZero(paths...)
	envPaths = slicesutil.Dedupe(envPaths)

	if len(envPaths) == 0 {
		return []string{}, nil
	} else if n < 1 {
		return envPaths, godotenv.Load(envPaths...)
	}
	loaded := []string{}
	for i, envPath := range envPaths {
		err := godotenv.Load(envPath)
		if err != nil {
			return loaded, err
		}
		loaded = append(loaded, envPath)
		if i+1 == n {
			break
		}
	}

	return loaded, nil
}

// GetDotEnvVal retrieves a single var from a `.env` file path
func GetDotEnvVal(envPath, varName string) (string, error) {
	cmd := fmt.Sprintf("grep %s '%s' | rev | cut -d= -f1 | rev", varName, envPath)

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %s", cmd)
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
		isFile, err := osutil.IsFile(fixedPath, true)
		if err != nil {
			return fixedPath, err
		} else if !isFile {
			return fixedPath, fmt.Errorf("path is not a file or is 0 size [%v]", fixedPath)
		}
		return fixedPath, nil
	}

	envPath = strings.TrimSpace(envPath)
	if len(envPath) > 0 {
		isFile, err := osutil.IsFile(fixedPath, true)
		if err != nil {
			return envPath, err
		} else if !isFile {
			return envPath, fmt.Errorf("path is not a file or is 0 size [%v]", envPath)
		}
		return envPath, nil
	}

	thisDirPath := "./.env"
	isFile, err := osutil.IsFile(thisDirPath, true)
	if err != nil {
		return thisDirPath, err
	} else if !isFile {
		return thisDirPath, fmt.Errorf("path is not a file or is 0 size [%v]", thisDirPath)
	}
	return thisDirPath, nil
}

type DotEnvOpts interface {
	DotEnvFilename() string
}

func ParseFlagsAndLoadDotEnv(opts DotEnvOpts) error {
	_, err := flags.Parse(opts)
	if err != nil {
		return err
	}
	if envFilename := opts.DotEnvFilename(); envFilename != "" {
		_, err := LoadDotEnv([]string{envFilename}, 1)
		return err
	}
	return nil
}
