package reader

import "github.com/kevinburke/ssh_config"

// NewConfigReader returns a new Reader of an ssh_config.Config
func NewConfigReader(config *ssh_config.Config) Reader {
	return config
}

// NewDefaultReader returns a new reader which contains default settings
func NewDefaultReader() Reader {
	return NewUserSettingsReader(ssh_config.DefaultUserSettings)
}

// NewUserSettingsReader returns a new reader of a ssh_config.UserSettings
func NewUserSettingsReader(settings *ssh_config.UserSettings) Reader {
	return userSettingsReader{settings: settings}
}

// userSettingsReader implements Reader
type userSettingsReader struct{ settings *ssh_config.UserSettings }

func (settings userSettingsReader) Get(alias, key string) (string, error) {
	return settings.settings.GetStrict(alias, key)
}

func (settings userSettingsReader) GetAll(alias, key string) ([]string, error) {
	return settings.settings.GetAllStrict(alias, key)
}
