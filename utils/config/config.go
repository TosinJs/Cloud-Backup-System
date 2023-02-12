package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	PORT       string `validate:"required"`
	DSN        string `validate:"required"`
	JWTSECRET  string `validate:"required"`
	AWS_ID     string `validate:"required"`
	AWS_SECRET string `validate:"required"`
	AWS_REGION string `validate:"required"`
	AWS_TOKEN  string
}

func LoadConfig(path, name, configType string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	var C Config

	err := viper.ReadInConfig()
	if err != nil {
		return C, err
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		return C, err
	}

	validate := validator.New()
	err = validate.Struct(&C)

	return C, err
}
