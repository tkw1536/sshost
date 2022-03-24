package source

import "github.com/tkw1536/stringreader"

// Source represents a source of ssh configuration settings
type Source interface {
	Alias(alias string) stringreader.Source
	isSSHSource()
}
