// Package settings описывает настройки для приложения.
package settings

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	defaultConfigFile = "config.yaml"
)

// Settings описывает структуру для хранения настроек сервера.
type Settings struct {
	API     APISettings          `koanf:"api"`
	Storage LocalStorageSettings `koanf:"localstorage"`
	Log     LogSettings          `koanf:"log"`
}

// APISettings подструктура для хранения настроек API.
type APISettings struct {
	Address string `koanf:"address"`
	Port    int    `koanf:"port"`
}

// LocalStorageSettings подструктура для хранения пути к локальному хранилищу.
type LocalStorageSettings struct {
	Path string `koanf:"path"`
}

// LogSettings подструктура для хранения настроек логгера.
type LogSettings struct {
	Level   string `koanf:"level"`
	Verbose bool   `koanf:"verbose"`
	Format  string `koanf:"format"`
}

// NewSettings принимает путь до файла настроек и пытается создать объект Settings.
func NewSettings(config string) (*Settings, error) {
	if config == "" {
		config = defaultConfigFile
	}

	k := koanf.New(".")

	// Read configuration file if exists.
	err := k.Load(file.Provider(config), yaml.Parser())
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", config, err)
	}

	s := &Settings{}

	err = k.Unmarshal("", &s)
	if err != nil {
		return nil, fmt.Errorf("unmarshal configuration: %w", err)
	}

	return s, nil
}
