package sshost

import (
	"fmt"

	"github.com/tkw1536/stringreader"
)

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
	"IPQoS",
	"KbdInteractiveDevices",
	"KnownHostsCommand",
	"LocalCommand",
	"LocalForward",
	// "LogLevel", // TODO: Can we safely ignore this?
	"PermitRemoteOpen",
	"PKCS11Provider",
	// "PreferredAuthentications", // TODO: Support authentications properly!
	"ProxyCommand",
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
	"NoHostAuthenticationForLocalhost",
	"PermitLocalCommand",
	"StreamLocalBindUnlink",
	"Tunnel", // TODO: Must be no!
	// "TunnelDevice",
	"UpdateHostKeys", // TODO: May have other values, but must be "no"
	"VisualHostKey",
}

// checkUnsupoportedConfig checks if source contains any unsupported configuration values.
// When this is the case, returns an error of type ErrUnsupportedConfig, else returns nil.
func checkUnsupportedConfig(source stringreader.Source) error {
	// check for unsupported flags (options that must be "no")
	for _, setting := range unsupportedFlags {
		value, ok := source.Lookup(setting)
		if ok && value == "yes" {
			return ErrUnsupportedConfig{Setting: setting, Value: "yes", Specific: true}
		}
	}

	// check for unsupported configs
	for _, setting := range unsupportedConfigs {
		value, ok := source.Lookup(setting)
		if ok && value != "" {
			return ErrUnsupportedConfig{Setting: setting, Value: value, Specific: false}
		}
	}

	return nil
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
