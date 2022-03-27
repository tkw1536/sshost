// Package sshost provides facilities for reading connecting to ssh hosts defined in ssh_config settings
package sshost

import (
	"os"
	"os/user"

	"github.com/kevinburke/ssh_config"
	"github.com/tkw1536/sshost/internal/pkg/source"
)

// NewDefaultEnvironment creates a new environment instance from the runtime environment.
//
// It reads both ~/.ssh and /etc/ssh as a source, see ssh_config.DefaultUserSettings.
// It uses operating system environment for defaults.
func NewDefaultEnvironment() (*Environment, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &Environment{
		Source: source.FromUserSettings(ssh_config.DefaultUserSettings),
		Strict: true,
		Defaults: Defaults{
			Username: user.Username,
		},
		Variables: os.Getenv,
	}, nil
}
