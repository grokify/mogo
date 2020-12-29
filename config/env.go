package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env"
	iom "github.com/grokify/simplego/io/ioutilmore"
	"github.com/joho/godotenv"
)

// EnvFileToJSONFile Converts an .env file to a JSON file using the definition
// provided in data.
func EnvFileToJSONFile(data interface{}, filepathENV, filepathJSON string, perm os.FileMode, prefix, indent string) error {
	err := godotenv.Load(filepathENV)
	if err != nil {
		return err
	}

	err = env.Parse(data)
	if err != nil {
		return err
	}

	return iom.WriteFileJSON(filepathJSON, data, perm, prefix, indent)
}

// Return a merged environment var which is split into multiple
// vars. This is useful when the system has a size limit on
// environment variables.
func JoinEnvNumbered(prefix, delimiter string, startInt uint8, includeBase bool) string {
	vals := []string{}
	if includeBase {
		val := os.Getenv(prefix)
		if len(val) > 0 {
			vals = append(vals, val)
		}
	}
	i := startInt
	for {
		val := os.Getenv(fmt.Sprintf("%s_%d", prefix, i))
		if len(val) > 0 {
			vals = append(vals, val)
		} else {
			break
		}
		i++
	}
	return strings.Join(vals, delimiter)
}
