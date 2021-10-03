package sshost

import (
	"fmt"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Configuration represents a connection to a single host
type Profile struct {
	env *Context

	// configuration, accessed only with GetConfig()
	config      Config
	configError error
	configValid sync.Once
}

// SetConfig sets the configuration for this Profile
func (profile *Profile) SetConfig(config Config) {
	profile.configValid = sync.Once{}
	profile.config = config
	profile.configError = nil
}

// GetConfig gets the configuration for the profile, and calls validate when needed
func (profile *Profile) GetConfig() (Config, error) {
	// run the validation once!
	profile.configValid.Do(func() {
		profile.configError = profile.config.Validate(profile.env.Strict)
	})

	// check if the validation was ok
	if profile.configError != nil {
		return Config{}, profile.configError
	}
	return profile.config, nil
}

// Dial creates a new net.Connection to host representing this profile
func (profile *Profile) Dial() (net.Conn, *ClosableStack, error) {
	cfg, err := profile.GetConfig()
	if err != nil {
		return nil, nil, err
	}

	network := cfg.AddressFamily.Network()
	if network == "" {
		return nil, nil, ErrUnknownAddressFamily
	}

	conn, err := net.DialTimeout(network, fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), cfg.ConnectTimeout)
	if err != nil {
		return nil, nil, err
	}

	return conn, NewClosableStack(conn), nil
}

// Config creates a new ssh configuration to use for a connection
func (profile Profile) Config() (*ssh.ClientConfig, error) {
	cfg, err := profile.GetConfig()
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: cfg.Username,

		Timeout: cfg.ConnectTimeout,

		// TODO: Implement security
		HostKeyAlgorithms: cfg.HostKeyAlgorithms,
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),

		Config: ssh.Config{
			Ciphers:      cfg.Ciphers,
			KeyExchanges: cfg.KexAlgorithms,
			MACs:         cfg.MACs,
		},
	}

	// when configure, setup a connection for an identity agent!
	identityAgent := cfg.IdentityAgent
	switch {
	case identityAgent != "" && identityAgent[0] == '$':
		identityAgent = profile.env.getenv(identityAgent[1:])
		fallthrough
	case identityAgent != "none":
		agentc, err := net.Dial("unix", identityAgent)
		if err != nil {
			return nil, err
		}
		client := agent.NewClient(agentc)
		config.Auth = append(config.Auth, ssh.PublicKeysCallback(client.Signers))
	}

	return config, nil
}

// Connect connects to the provided host using the given connection.
func (profile Profile) Connect(conn net.Conn) (*ssh.Session, error) {
	config, err := profile.Config()
	if err != nil {
		return nil, err
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, conn.RemoteAddr().String(), config)
	if err != nil {
		return nil, err
	}

	client := ssh.NewClient(c, chans, reqs)
	return client.NewSession()
}
