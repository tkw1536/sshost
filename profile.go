package sshost

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/tkw1536/sshost/pkg/closer"

	"golang.org/x/crypto/ssh"
)

// Profile represents a connection to a single host
type Profile struct {
	env *Environment

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

var ErrContextClosed = errors.New("Profile.Dial: Context was closed")

// Dial creates a new net.Conn to the host behind the given profile.
// When the profile contains a JumpHost, this might involve connecting to other ssh hosts.
// If the context is cancelled, the connection to any existing ssh host is closed.
//
// Proxy indiciates an ssh proxy to dial the connection from.
// When proxy is nil, does not use a proxy.
func (profile *Profile) Dial(proxy *ssh.Client, ctx context.Context) (net.Conn, *closer.Stack, error) {
	// shortcut: if the context is already closed, bail out immediatly!
	if ctx.Err() != nil {
		return nil, nil, ErrContextClosed
	}

	// create a stack and current connection
	// used to establish the connection and register all the closers!
	stack := closer.NewStack()
	hop := proxy

	// keep track of the current connection attempt
	// and close the stack at the end of the connection!
	doneC := make(chan struct{})
	defer close(doneC)
	var cancelDial uint32 // non-zero when cancelled
	go func() {
		select {
		case <-ctx.Done():
			atomic.StoreUint32(&cancelDial, 1)
			stack.Close()
		case <-doneC: /* connection established */
		}
	}()

	// iterate over all the hops
	var err error
	var jumpStack *closer.Stack
	for _, jumpHost := range profile.config.ProxyJump {
		if atomic.LoadUint32(&cancelDial) == 0 {
			hop, jumpStack, err = profile.env.NewClient(hop, jumpHost, ctx)
			stack.PushStack(jumpStack)
		} else {
			err = ErrContextClosed
		}

		if err != nil {
			defer stack.Close()
			return nil, nil, err
		}
	}

	// determine the parameters for the final "real" hop
	cfg, err := profile.GetConfig()
	if err != nil {
		return nil, nil, err
	}

	network := cfg.AddressFamily.Network()
	if network == "" {
		return nil, nil, ErrUnknownAddressFamily
	}
	address := fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)

	// establish the connection from the final hop to the machine itself
	// do this either via the real network, or via the existing client
	var conn net.Conn
	if atomic.LoadUint32(&cancelDial) == 0 {
		if hop == nil {
			conn, err = net.DialTimeout(network, address, cfg.ConnectTimeout)
		} else {
			conn, err = hop.Dial(network, address)
		}
	} else {
		err = ErrContextClosed
	}

	if err != nil {
		defer stack.Close()
		return nil, nil, err
	}

	stack.Push(conn)
	return conn, stack, nil
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

		Auth: profile.env.Auth.Methods(cfg.PreferredAuthentications, profile),
	}

	return config, nil
}

// Connect connects to the provided host using the given connection.
func (profile Profile) Connect(conn net.Conn) (*ssh.Client, error) {
	config, err := profile.Config()
	if err != nil {
		return nil, err
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, conn.RemoteAddr().String(), config)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, chans, reqs), nil
}
