package sshost

import (
	"fmt"
	"os"
	"os/user"

	"github.com/tkw1536/sshost/reader"
)

// Context represents a context to derive a Config from
type Context struct {
	Reader reader.Reader

	Strict bool

	DefaultUsername string
	Variables       func(name string) string
}

func NewDefaultContext() (*Context, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &Context{
		Reader:          reader.NewDefaultReader(),
		DefaultUsername: user.Username,
		Variables:       os.Getenv,
	}, nil
}

// getenv returns ctx.Variables, protected against Variables being nil
func (ctx *Context) getenv(name string) string {
	if ctx.Variables == nil {
		return ""
	}
	return ctx.Variables(name)
}

// list of security-critical unsupported configs
var unsupportedConfigs = []string{
	// "AddKeysToAgent", // we don't need to retain access to keys
	// "BatchMode", // always in batch mode, connection may fail
	"BindAddress",
	"CanonicalDomains",
	// "CanonicalizeFallbackLocal",
	// "CanonicalizeHostname",
	// "CanonicalizeMaxDots",
	"CanonicalizePermittedCNAMEs",
	"CASignatureAlgorithms",
	"CertificateFile",
	// "CheckHostIP", // TODO: implement me!
	// "ClearAllForwardings",

	"DynamicForward",
	// "EscapeChar", // TODO: Support by user
	// "FingerprintHash", // TODO: Just used for output!

	// "GlobalKnownHostsFile", // TODO: Support me!
	// "HostbasedAcceptedAlgorithms",
	"HostKeyAlias",
	// "IdentityFile", // TODO: Support me!
	"IPQoS",
	// "KbdInteractiveAuthentication", // TODO: Support me!
	"KbdInteractiveDevices",
	"KnownHostsCommand",
	"LocalCommand",
	"LocalForward",
	// "LogLevel", // TODO: Can we safely ignore this?
	// "NumberOfPasswordPrompts", // TODO: Implement me!
	// "PasswordAuthentication", // TODO: Implement me!
	"PermitRemoteOpen",
	"PKCS11Provider",
	// "PreferredAuthentications", // TODO: Support authentications properly!
	"ProxyCommand",
	"ProxyJump",
	// "ProxyUseFdpass", // ProxyCommand is unsupported!
	"PubkeyAcceptedAlgorithms",
	// "PubkeyAuthentication", // TODO: Support authentication properly!
	// "RekeyLimit", // TODO: Support this properly!
	"RemoteCommand",
	"RemoteForward",
	"RequestTTY",
	"SendEnv",
	// "ServerAliveCountMax", // TODO: Support ServerAliveInterval
	"SessionType",
	"SetEnv",
	"StdinNull",
	// "StreamLocalBindMask",
	// "StreamLocalBindUnlink",
	// "StrictHostKeyChecking", // TODO: Support me!
	// "TCPKeepAlive", // TODO: Enabled by default!
	// "UserKnownHostsFile",  // TODO: Support authentication properly!
	// "VerifyHostKeyDNS", // TODO: Support properly!
	// "XAuthLocation", // TODO: Support authentication properly!
}

var unsupportedFlags = []string{
	"ControlMaster",
	// "ControlPath",
	// "ControlPersist",
	"ExitOnForwardFailure",
	"ForkAfterAuthentication",
	"ForwardAgent",
	"ForwardX11",
	// "ForwardX11Timeout",
	"ForwardX11Trusted",
	"GatewayPorts",
	"GSSAPIAuthentication",
	"GSSAPIDelegateCredentials",
	"HashKnownHosts",
	"HostbasedAuthentication",
	"IdentitiesOnly",
	"NoHostAuthenticationForLocalhost",
	"PermitLocalCommand",
	"StreamLocalBindUnlink",
	"Tunnel", // TODO: Must be no!
	// "TunnelDevice",
	"UpdateHostKeys", // TODO: May have other values, but must be "no"
	"VisualHostKey",
}

// NewProfile gets a new profile for the environment
func (ctx *Context) NewProfile(alias string) (profile *Profile, err error) {
	cfg, err := ctx.NewConfig(alias)
	if err != nil {
		return nil, err
	}
	return &Profile{
		env:    ctx,
		config: cfg,
	}, nil
}

// NewConfig reads a new configuration for the specific alias from the configuration
func (ctx Context) NewConfig(alias string) (cfg Config, err error) {
	cHostname, err := reader.Default(ctx.Reader, alias, "Hostname", alias)
	if err != nil {
		return cfg, err
	}

	cPort, err := reader.Uint(ctx.Reader, alias, "Port", 10, 16, 22)
	if err != nil {
		return cfg, err
	}

	cUsername, err := reader.Default(ctx.Reader, alias, "User", ctx.DefaultUsername)
	if err != nil {
		return cfg, err
	}

	cIdentityAgent, err := reader.Default(ctx.Reader, alias, "IdentityAgent", "$SSH_AUTH_SOCK")
	if err != nil {
		return cfg, err
	}

	cAddressFamily, err := reader.Default(ctx.Reader, alias, "AddressFamily", string(DefaultAddressFamily))
	if err != nil {
		return cfg, err
	}

	cCiphers, err := reader.StringSlice(ctx.Reader, alias, "Ciphers", nil)
	if err != nil {
		return cfg, err
	}

	cHostKeyAlgorithms, err := reader.StringSlice(ctx.Reader, alias, "HostKeyAlgorithms", nil)
	if err != nil {
		return cfg, err
	}

	cCompression, err := reader.YesNo(ctx.Reader, alias, "Compression", false)
	if err != nil {
		return cfg, err
	}

	cConnectionAttempts, err := reader.Uint(ctx.Reader, alias, "ConnectionAttempts", 10, 64, 1)
	if err != nil {
		return cfg, err
	}

	cConnectTimeout, err := reader.Seconds(ctx.Reader, alias, "ConnectTimeout", 0)
	if err != nil {
		return cfg, err
	}

	cKexAlgorithms, err := reader.StringSlice(ctx.Reader, alias, "KexAlgorithms", nil)
	if err != nil {
		return cfg, err
	}

	cMACs, err := reader.StringSlice(ctx.Reader, alias, "MACs", nil)
	if err != nil {
		return cfg, err
	}

	cServerAliveInterval, err := reader.Seconds(ctx.Reader, alias, "ServerAliveInterval", 0)
	if err != nil {
		return cfg, err
	}

	cServerAliveCountMax, err := reader.Uint(ctx.Reader, alias, "ServerAliveCountMax", 10, 64, 3)
	if err != nil {
		return cfg, err
	}

	// TODO: Support cRekeyLimit properly!
	cRekeyLimit, err := reader.Default(ctx.Reader, alias, "RekeyLimit", "default none")
	if err != nil {
		return cfg, err
	}

	if cRekeyLimit != "default none" {
		return cfg, ErrUnsupportedConfig{Setting: "RekeyLimit", Value: cRekeyLimit, Specific: true}
	}

	// check for unsupported flags (options that must be "no")
	for _, setting := range unsupportedFlags {
		value, err := reader.YesNo(ctx.Reader, alias, setting, false)
		if err != nil {
			return cfg, err
		}
		if value {
			return cfg, ErrUnsupportedConfig{Setting: setting, Value: "yes", Specific: true}
		}
	}

	// check for unsupported configs
	for _, setting := range unsupportedConfigs {
		value, err := ctx.Reader.Get(alias, setting)
		if err != nil {
			return cfg, err
		}
		if value != "" {
			return cfg, ErrUnsupportedConfig{Setting: setting, Value: value, Specific: false}
		}
	}

	return Config{
		AddressFamily:       AddressFamily(cAddressFamily),
		Ciphers:             cCiphers,
		Compression:         cCompression,
		ConnectionAttempts:  cConnectionAttempts,
		ConnectTimeout:      cConnectTimeout,
		HostKeyAlgorithms:   cHostKeyAlgorithms,
		Hostname:            cHostname,
		IdentityAgent:       cIdentityAgent,
		KexAlgorithms:       cKexAlgorithms,
		MACs:                cMACs,
		Port:                uint32(cPort),
		RekeyLimit:          cRekeyLimit,
		ServerAliveCountMax: cServerAliveCountMax,
		ServerAliveInterval: cServerAliveInterval,
		Username:            cUsername,
	}, nil
}

// ErrUnsupportedConfig represents an unsupported configuration setting
type ErrUnsupportedConfig struct {
	// Name of the setting that is unsupported
	Setting string

	// Value the unsupported setting currently has
	Value string

	// When true indicates that only this specific value is unsupported
	Specific bool
}

func (u ErrUnsupportedConfig) Error() string {
	if u.Specific {
		return fmt.Sprintf("unsupported configuration value for setting %q: %q", u.Setting, u.Value)
	}
	return fmt.Sprintf("unsupported configuration setting %q (has value %q)", u.Setting, u.Value)
}
