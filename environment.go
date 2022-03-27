package sshost

import (
	"context"

	"github.com/tkw1536/sshost/internal/pkg/closer"
	"github.com/tkw1536/sshost/internal/pkg/host"
	"github.com/tkw1536/sshost/internal/pkg/source"
	"golang.org/x/crypto/ssh"
)

// Environment represents a system environment to derive a Config from
type Environment struct {
	// Source is a source of configuration values
	Source source.Source

	// Strict is used to enable strict validation of settings.
	Strict bool

	// Defaults are the defaults for creating new profiles
	Defaults Defaults
	Auth     AuthEnv

	// Variables contains values of system environment variables
	Variables func(name string) string
}

// getenv returns ctx.Variables, protected against Variables being nil
func (env Environment) getenv(name string) string {
	if env.Variables == nil {
		return ""
	}
	return env.Variables(name)
}

// NewClient creates a new client.
// See also DialContext and connect.
//
// The provided context is only used during the dialing phase, if the context is canceled after the context phase, it has no effect.
func (env Environment) NewClient(proxy *ssh.Client, alias string, ctx context.Context) (*ssh.Client, *closer.Stack, error) {
	profile, err := env.NewProfile(alias)
	if err != nil {
		return nil, nil, err
	}

	conn, closers, err := profile.Dial(proxy, ctx)
	if err != nil {
		return nil, nil, err
	}

	client, err := profile.Connect(conn)
	if err != nil {
		defer closers.Close()
		return nil, nil, err
	}

	closers.Push(client)
	return client, closers, err
}

// NewProfile gets a new profile for the environment
func (env *Environment) NewProfile(alias string) (profile *Profile, err error) {
	cfg, err := env.NewConfig(alias)
	if err != nil {
		return nil, err
	}
	return &Profile{
		env:    env,
		config: cfg,
	}, nil
}

// NewConfig creates a new configuration for the provided alias.
//
// alias may be a simple hostname or a more complex ssh uri.
// See config.ParseHost for details.
func (env Environment) NewConfig(alias string) (Config, error) {
	// Parse the hostname
	h, err := host.ParseHost(alias)
	if err != nil {
		return Config{}, err
	}

	// create a new configuration
	src := env.Source.Alias(h.Host)
	cfg, err := NewConfig(src, h, env.Defaults)
	if err != nil {
		return cfg, err
	}

	// TODO: Expand configuration!

	return cfg, nil
}
