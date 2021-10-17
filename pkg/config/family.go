package config

import "errors"

// AddressFamily specifies which address family to use when connecting.
type AddressFamily string

const (
	DefaultAddressFamily AddressFamily = "any"
	IPv4AddressFamily    AddressFamily = "inet"
	IPv6AddressFamily    AddressFamily = "inet6"
)

// Valid checks if the provided AddressFamily is valid
func (a AddressFamily) Valid() bool {
	return a == DefaultAddressFamily || a == IPv4AddressFamily || a == IPv6AddressFamily
}

// Network returns the network corresponding to this AddressFamily.
// In case of an unknown AddressFamily, returns "".
func (a AddressFamily) Network() string {
	switch a {
	case "", DefaultAddressFamily:
		return "tcp"
	case IPv4AddressFamily:
		return "tcp4"
	case IPv6AddressFamily:
		return "tcp6"
	default:
		return ""
	}
}

// ErrUnknownAddressFamily represents an unknown AddressFamily
var ErrUnknownAddressFamily = errors.New("unknown AddressFamily")
