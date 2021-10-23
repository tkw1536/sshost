package config

import (
	"github.com/kevinburke/ssh_config"
	"github.com/tkw1536/stringreader"
)

// NewConfigSource returns a new source for the provided alias and configuration
func NewConfigSource(config *ssh_config.Config, alias string) stringreader.Source {
	return configSource{alias: alias, config: config}
}

type configSource struct {
	config *ssh_config.Config
	alias  string
}

func (src configSource) Get(key string) (value string, ok bool) {
	value, err := src.config.Get(src.alias, key)
	if err != nil {
		return "", false
	}
	return value, true
}

func (src configSource) GetAll(key string) (value []string, ok bool) {
	value, err := src.config.GetAll(src.alias, key)
	if err != nil {
		return nil, false
	}
	return value, true
}

// NewUserSettingsSource returns a new source for the provided alias in the user settings
func NewUserSettingsSource(settings *ssh_config.UserSettings, alias string) stringreader.Source {
	return userSettingsSource{settings: settings, alias: alias}
}

type userSettingsSource struct {
	settings *ssh_config.UserSettings
	alias    string
}

func (src userSettingsSource) Get(key string) (value string, ok bool) {
	value, err := src.settings.GetStrict(src.alias, key)
	if err != nil {
		return "", false
	}
	return value, true
}

func (src userSettingsSource) GetAll(key string) (value []string, ok bool) {
	value, err := src.settings.GetAllStrict(src.alias, key)
	if err != nil {
		return nil, false
	}
	return value, true
}
