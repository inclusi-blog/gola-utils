package configuration_loader

import (
	"encoding/json"
	"github.com/gola-glitch/gola-utils/logging"
	"os"

	"gopkg.in/go-playground/validator.v9"
)

type ConfigLoader interface {
	Load(filePath string, configStruct interface{}) error
}

type configLoader struct {
}

func NewConfigLoader() ConfigLoader {
	return configLoader{}
}

func (configLoader configLoader) Load(filePath string, configStruct interface{}) error {
	logger := logging.NewLoggerEntry()
	file, fileOpeningErr := os.Open(filePath)
	if fileOpeningErr != nil {
		logger.Error("Err while opening config file : ", fileOpeningErr)
		return fileOpeningErr
	}

	decoder := json.NewDecoder(file)
	decodingErr := decoder.Decode(&configStruct)
	if decodingErr != nil {
		logger.Error("Err while decoding config file : ", decodingErr)
		return decodingErr
	}

	validate := validator.New()
	validationErr := validate.Struct(configStruct)
	if validationErr != nil {
		logger.Error("Err while validating config file : ", validationErr)
		return validationErr
	}
	return validationErr
}
