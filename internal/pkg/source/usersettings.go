package source

import (
	"github.com/kevinburke/ssh_config"
	"github.com/tkw1536/stringreader"
)

// FromUserSettings creates a new source for the provided UserSettings and alias
func FromUserSettings(settings *ssh_config.UserSettings, alias string) stringreader.Source {
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
