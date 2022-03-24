// Package source provides a mapping between an ssh files and config sources
package source

import (
	"github.com/kevinburke/ssh_config"
	"github.com/tkw1536/stringreader"
)

func FromSSHConfig(config *ssh_config.Config) Source {
	return sshConfig{config: config}
}

type sshConfig struct {
	config *ssh_config.Config

	aliasSet bool
	alias    string
}

// isSSHConfig indicates that sshConfig is an ssh source
func (sshConfig) isSSHSource() {}

func (config sshConfig) Alias(alias string) stringreader.Source {
	config.alias = alias
	config.aliasSet = true
	return config
}

func (config sshConfig) Get(key string) (value string, ok bool) {
	if !config.aliasSet {
		return "", false
	}

	value, err := config.config.Get(config.alias, key)
	if err != nil {
		return "", false
	}
	return value, true
}

func (config sshConfig) GetAll(key string) (value []string, ok bool) {
	if !config.aliasSet {
		return nil, false
	}

	value, err := config.config.GetAll(config.alias, key)
	if err != nil {
		return nil, false
	}
	return value, true
}
