package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type CliConfig struct {
	Twitter struct {
		Cookie        string `mapstructure:"cookie"`
		CsrfToken     string `mapstructure:"X_CSRF_TOKEN"`
		Authorization string `mapstructure:"Authorization"`
	} `mapstructure:"twitter"`
	Pixiv struct {
		Cookie string `mapstructure:"PIXIV_COOKIE"`
		Agent  string `mapstructure:"PIXIV_USER_AGENT"`
	} `mapstructure:"pixiv"`
}

var Config CliConfig

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$XDG_CONFIG_HOME/wormface-cli")
	viper.AddConfigPath("$HOME/.config/wormface-cli")
	viper.AddConfigPath("%APPDATA%/wormface-cli")

	// write default
	if err := viper.ReadInConfig(); err != nil {
		var notFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &notFoundErr) {
			return err
		}

		var defaultConfigHome string
		if runtime.GOOS == "windows" {
			defaultConfigHome = filepath.Join(os.Getenv("APPDATA"), "wormface-cli")
		} else {
			defaultConfigHome = filepath.Join(os.Getenv("HOME"), ".config", "wormface-cli")
		}
		var defaultConfigPath = filepath.Join(defaultConfigHome, "config.yaml")

		if err := os.MkdirAll(defaultConfigHome, 0755); err != nil {
			return err
		}

		if _, err := os.Create(defaultConfigPath); err != nil {
			return err
		}

		if err := viper.WriteConfigAs(defaultConfigPath); err != nil {
			return err
		}
	}

	return viper.Unmarshal(&Config)
}
