// Package config contains Config and related datastructures.
//
// All types in this configuration are considered stateless.
package config

import (
	"errors"
	"fmt"
	"time"
)

// Config represents a configuration for a single host
type Config struct {
	AddressFamily       AddressFamily
	Ciphers             []string
	Compression         bool
	ConnectionAttempts  uint64
	ConnectTimeout      time.Duration
	HostKeyAlgorithms   []string
	Hostname            string
	IdentityAgent       string
	KexAlgorithms       []string
	MACs                []string
	ProxyJump           []string
	Port                uint16
	RekeyLimit          string // TODO: Proper datatype
	ServerAliveCountMax uint64
	ServerAliveInterval time.Duration
	Username            string
}

var knownExchangeAlgos = []string{
	"diffie-hellman-group1-sha1",
	"diffie-hellman-group14-sha1",
	"ecdh-sha2-nistp256",
	"ecdh-sha2-nistp384",
	"ecdh-sha2-nistp521",
	"curve25519-sha256@libssh.org",
	"diffie-hellman-group-exchange-sha1",
	"diffie-hellman-group-exchange-sha256",
}
var knownMACs = []string{
	"hmac-sha2-256-etm@openssh.com", "hmac-sha2-256", "hmac-sha1", "hmac-sha1-96",
}
var knownCiphers = []string{
	"aes128-ctr", "aes192-ctr", "aes256-ctr",
	"aes128-gcm@openssh.com",
	"chacha20-poly1305@openssh.com",
	"arcfour256", "arcfour128", "arcfour",
	"aes128-cbc",
	"3des-cbc",
}

// Validate validates the provided configuration and (where necessary) normalizes it.
// When validation fails, returns an error of type ErrField; otherwise err is nil.
//
// When strict is false, if no algorithms selected within the configuration are supported uses default algorithms instead.
// When strict is true, an error is returned instead.
func (cfg *Config) Validate(strict bool) (err error) {
	if !cfg.AddressFamily.Valid() {
		return NewErrField(nil, "AddressFamily")
	}
	if err := filterSliceField(&cfg.Ciphers, strict, "Ciphers", knownCiphers); err != nil {
		return err
	}
	if cfg.Compression {
		return NewErrField(nil, "Compression")
	}
	if cfg.ConnectionAttempts != 1 {
		return NewErrField(nil, "ConnectionAttempts")
	}
	// ConnectTimeout: no validation
	if err := filterSliceField(&cfg.HostKeyAlgorithms, strict, "HostKeyAlgorithms", knownExchangeAlgos); err != nil {
		return err
	}
	if cfg.Hostname == "" {
		return NewErrField(errEmptyField, "Hostname")
	}
	// IdentityAgent: no validation
	if err := filterSliceField(&cfg.KexAlgorithms, strict, "KexAlgorithms", knownExchangeAlgos); err != nil {
		return err
	}
	if err := filterSliceField(&cfg.MACs, strict, "MACs", knownMACs); err != nil {
		return err
	}
	for _, pj := range cfg.ProxyJump {
		if !ValidHost(pj) {
			return NewErrField(nil, "ProxyJump")
		}
	}
	if cfg.Port == 0 || cfg.Port >= 65535 {
		return NewErrField(nil, "Port")
	}
	if cfg.RekeyLimit != "default none" {
		return NewErrField(nil, "RekeyLimit")
	}
	// ServerAliveCountMax: no validation
	if cfg.ServerAliveInterval != 0 {
		return NewErrField(errEmptyField, "ServerAliveInterval")
	}
	if cfg.Username == "" {
		return NewErrField(errEmptyField, "Username")
	}
	return
}

// filterSliceField calls filterSlice, and ignores errors unless strict = True
func filterSliceField(slice *[]string, strict bool, field string, valid []string) error {
	var err error
	*slice, err = filterSlice(*slice, valid)
	if err != nil {
		if strict {
			return NewErrField(err, field)
		}
		*slice = nil
	}
	return nil
}

// filterSlice filters s by elements only in valid.
// Does not re-allocate, and invalidates memory used by s.
//
// When slice is nil, never returns an error.
// When slice becomes empty, returns errAllUnsupported
func filterSlice(slice []string, valid []string) ([]string, error) {
	// special case: 0-size slice or valid
	if len(slice) == 0 || len(valid) == 0 {
		return slice, nil
	}

	// cache which elements exist
	cache := make(map[string]struct{}, len(valid))
	for _, v := range valid {
		cache[v] = struct{}{}
	}

	// filter s according to the cache
	result := slice[:0]
	for _, element := range slice {
		if _, ok := cache[element]; ok {
			result = append(result, element)
		}
	}

	if len(result) == 0 {
		return slice, errAllUnsupported
	}

	return result, nil
}

// ErrField represents an error for the provided field
type ErrField struct {
	error
	Field string
}

// NewErrField creates a new ErrField for the given field and error.
// When err is nil, picks a generic error.
func NewErrField(err error, Field string) ErrField {
	if err == nil {
		err = errInvalidField
	}
	return ErrField{error: err, Field: Field}
}

var errInvalidField = errors.New("field value invalid")
var errEmptyField = errors.New("field must be non-empty")
var errAllUnsupported = errors.New("no value is supported")

func (err ErrField) Unwrap() error {
	return err.error
}

func (err ErrField) Error() string {
	return fmt.Sprintf("Field %q: %s", err.Field, err.error.Error())
}
