package apiserver

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"golang-auth/internal/app/repository"
	"golang-auth/internal/app/service"
	"io"
	"os"
)

type Config struct {
	Addr string            `json:"addr" validate:"required"`
	DB   repository.Config `json:"db" validate:"required"`
	JWT  service.JWTConfig `json:"jwt" validate:"required"`
}

func LoadConfig() (*Config, error) {
	configPath := "./configs/config.json"
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, err
	}

	return &config, nil

}
