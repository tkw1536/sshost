package source

import (
	"github.com/kevinburke/ssh_config"
	"github.com/tkw1536/stringreader"
)

func FromUserSettings(settings *ssh_config.UserSettings) Source {
	return sshUserSettings{settings: settings}
}

type sshUserSettings struct {
	settings *ssh_config.UserSettings

	aliasSet bool
	alias    string
}

// isSSHConfig indicates that sshUserSettings is an ssh source
func (sshUserSettings) isSSHSource() {}

func (settings sshUserSettings) Alias(alias string) stringreader.Source {
	settings.alias = alias
	settings.aliasSet = true
	return settings
}

func (settings sshUserSettings) Get(key string) (value string, ok bool) {
	if !settings.aliasSet {
		return "", false
	}

	value, err := settings.settings.GetStrict(settings.alias, key)
	if err != nil {
		return "", false
	}
	return value, true
}

func (settings sshUserSettings) GetAll(key string) (value []string, ok bool) {
	if !settings.aliasSet {
		return nil, false
	}

	value, err := settings.settings.GetAllStrict(settings.alias, key)
	if err != nil {
		return nil, false
	}
	return value, true
}
