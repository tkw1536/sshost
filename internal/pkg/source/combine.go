package source

import "github.com/tkw1536/stringreader"

type csource struct {
	sources []Source
	aliases []stringreader.Source
}

func (csource) isSSHSource() {}

func (c csource) Alias(alias string) stringreader.Source {
	c.aliases = make([]stringreader.Source, len(c.sources))
	for i, s := range c.sources {
		c.aliases[i] = s.Alias(alias)
	}
	return c
}

func (c csource) Lookup(key string) (value string, ok bool) {
	for _, a := range c.aliases {
		value, ok = a.Lookup(key)
		if ok {
			return
		}
	}
	return "", false
}
func (c csource) LookupAll(key string) (values []string, ok bool) {
	for _, a := range c.aliases {
		values, ok = a.LookupAll(key)
		if ok {
			return
		}
	}
	return nil, false
}
