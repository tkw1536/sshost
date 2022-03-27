package source

import "github.com/tkw1536/stringreader"

// Source is a source of configuration values.
//
// All Source types are defined in this package.
// New Sources can be created using NewSourceMap, FromSSHConfig and FromUserSettings.
// Sources can be combined using Combine.
type Source interface {
	// Alias returns source for a specific alias
	//
	// When an alias does not exist, should return default values.
	Alias(alias string) stringreader.Source
	isSSHSource()
}

// NewSource returns the same glo
func NewSourceMap(globals map[string]string) Source {
	return smap(globals)
}

type smap map[string]string

func (m smap) isSSHSource() {}

func (m smap) Alias(alias string) stringreader.Source {
	return stringreader.SourceSmartSplit{
		SourceSingle: stringreader.SourceSingleMap(m),
	}
}

func (m smap) Get(key string) (value string, ok bool) {
	value, ok = m[key]
	return
}
func (m smap) GetAll(key string) (values []string, ok bool) {
	value, ok := m[key]
	if !ok {
		return nil, true
	}
	return []string{value}, true
}
