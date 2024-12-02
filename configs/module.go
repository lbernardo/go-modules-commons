package configs

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"log"
	"os"
)

var Module = fx.Provide(
	registerAppConfigs,
)

func registerAppConfigs() *viper.Viper {
	envConfig := os.Getenv("APP_CONFIG_FILE")
	configsFiles := []string{
		envConfig,
		"config.dev.yaml",
		"/etc/secrets/config.yaml",
	}
	vp, err := readConfigByPossibleFiles(configsFiles)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return vp
}

func readConfigByPossibleFiles(possibleFiles []string) (*viper.Viper, error) {
	var filename string
	for _, file := range possibleFiles {
		if file == "" {
			continue
		}
		_, err := os.Stat(file)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		filename = file
	}
	if filename == "" {
		return nil, fmt.Errorf("config file not found")
	}

	vp := viper.New()

	vp.SetConfigFile(filename)
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error read config: %v", err)
	}
	return vp, nil
}

func ReadConfigFile(possibleFiles []string, result any) error {
	vp, err := readConfigByPossibleFiles(possibleFiles)
	if err != nil {
		return err
	}
	if err := vp.Unmarshal(&result); err != nil {
		return err
	}
	return nil
}
