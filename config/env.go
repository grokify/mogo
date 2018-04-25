package config

import (
	"os"

	"github.com/caarlos0/env"
	iom "github.com/grokify/gotilla/io/ioutilmore"
	"github.com/joho/godotenv"
)

// EnvFileToJSONFile Converts an .env file to a JSON file using the definition
// provided in data.
func EnvFileToJSONFile(data interface{}, filepathENV, filepathJSON string, perm os.FileMode, pretty bool) error {
	err := godotenv.Load(filepathENV)
	if err != nil {
		return err
	}

	err = env.Parse(data)
	if err != nil {
		return err
	}

	return iom.WriteJSON(filepathJSON, data, perm, pretty)
}
