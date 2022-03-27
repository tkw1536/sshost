package sshost

import (
	"errors"
	"fmt"

	"github.com/tkw1536/sshost/internal/pkg/host"
	"github.com/tkw1536/sshost/internal/pkg/slices"
)

// list of algorhtms supported for specific fields
var sKeyAlgorithms = slices.Combine(knownCertAlgos, knownKeyAlgos, knownKexAlgos)
var sKexAlgorithms = slices.Combine(knownKexAlgos)
var sCiphers = slices.Combine(knownCiperNames)
var sMACs = slices.Combine(knownMACNames)

// Validate validates the provided configuration and normalizes it.
// When validation fails, returns an error of type ErrField; otherwise err is nil.
//
// When strict is false, if no algorithms selected within the configuration are supported uses default algorithms instead.
// When strict is true, an error is returned instead.
func (cfg *Config) Validate(strict bool) (err error) {
	if !cfg.AddressFamily.Valid() {
		return NewErrField(nil, "AddressFamily")
	}
	if err := filterSliceField(&cfg.Ciphers, strict, "Ciphers", sCiphers); err != nil {
		return err
	}
	if cfg.Compression {
		return NewErrField(nil, "Compression")
	}
	if cfg.ConnectionAttempts != 1 {
		return NewErrField(nil, "ConnectionAttempts")
	}
	// ConnectTimeout: no validation
	if err := filterSliceField(&cfg.HostKeyAlgorithms, strict, "HostKeyAlgorithms", sKeyAlgorithms); err != nil {
		return err
	}
	if cfg.Hostname == "" {
		return NewErrField(errEmptyField, "Hostname")
	}
	// IdentityAgent: no validation
	if err := filterSliceField(&cfg.KexAlgorithms, strict, "KexAlgorithms", knownKexAlgos); err != nil {
		return err
	}
	if err := filterSliceField(&cfg.MACs, strict, "MACs", sMACs); err != nil {
		return err
	}
	for _, pj := range cfg.ProxyJump {
		if !host.ValidHost(pj) {
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
	*slice, err = slices.Filter(*slice, valid)
	if err != nil {
		if strict {
			return NewErrField(err, field)
		}
		*slice = nil
	}
	return nil
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

func (err ErrField) Unwrap() error {
	return err.error
}

func (err ErrField) Error() string {
	return fmt.Sprintf("Field %q: %s", err.Field, err.error.Error())
}
