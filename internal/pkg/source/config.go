// Package source provides a mapping between an ssh files and config sources
package source

import (
	"github.com/kevinburke/ssh_config"
	"github.com/tkw1536/stringreader"
)

// FromConfig creates a new source for the provided ssh_config and alias
func FromConfig(config *ssh_config.Config, alias string) stringreader.Source {
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
