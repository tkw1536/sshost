package sshost

import (
	"context"
	"os"
	"os/user"

	"github.com/kevinburke/ssh_config"
	"github.com/tkw1536/sshost/pkg/closer"
	"github.com/tkw1536/sshost/pkg/host"
	"github.com/tkw1536/stringreader"
	"golang.org/x/crypto/ssh"
)

// Environment represents a system environment to derive a Config from
type Environment struct {
	// Settings and Config are used to read data from the environment.
	// Exactly one of Settings and Config must be non-nil.
	//
	// If both are nil, any function of the environment may panic().
	// If both are non-nil, the behavior is undefined.
	Settings *ssh_config.UserSettings
	Config   *ssh_config.Config

	// Strict is used to enable strict validation of settings.
	Strict bool

	// Defaults are the defaults for creating new profiles
	Defaults Defaults
	Auth     AuthEnv

	// Variables contains values of system environment variables
	Variables func(name string) string
}

func NewDefaultEnvironment() (*Environment, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &Environment{
		Settings: ssh_config.DefaultUserSettings,
		Strict:   false,
		Defaults: Defaults{
			Username: user.Username,
		},
		Variables: os.Getenv,
	}, nil
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
func (env *Environment) NewClient(proxy *ssh.Client, alias string, ctx context.Context) (*ssh.Client, *closer.Stack, error) {
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

func (env Environment) source(alias string) stringreader.Source {
	if env.Config != nil {
		return NewConfigSource(env.Config, alias)
	}
	if env.Settings != nil {
		return NewUserSettingsSource(env.Settings, alias)
	}

	panic("env.Source: env.settings and env.config are nil")
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
	src := env.source(h.Host)
	cfg, err := NewConfig(src, h, env.Defaults)
	if err != nil {
		return cfg, err
	}

	// TODO: Expand configuration!

	return cfg, nil
}
