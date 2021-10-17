package config

import (
	"errors"
	"strconv"
	"strings"
)

// Host represents a host connection profile.
// It is parsed from a simple URL, see ParseHost.
type Host struct {
	User string
	Host string
	Port uint16

	// TODO: Additional configuration settings!
}

var ErrHostUnsupported = errors.New("unsupported hostname")

// ParseHost parses a host from a string
func ParseHost(host string) (h Host, err error) {
	if strings.Contains(host, "://") {
		return h, ErrHostUnsupported
	}
	h.Host = host

	// trim off the '@' sign
	index := strings.IndexRune(h.Host, '@')
	if index >= 0 {
		h.User = h.Host[:index]
		h.Host = h.Host[index+1:]
	}

	index = strings.IndexRune(h.Host, ':')
	if index >= 0 {
		port := h.Host[index+1:]
		lport, err := strconv.ParseUint(port, 10, 16)
		if err != nil {
			return Host{}, err
		}
		h.Port = uint16(lport)

		h.Host = h.Host[:index]
	}

	return
}

// ValidHost checks if Host is a valid host
func ValidHost(host string) bool {
	h, err := ParseHost(host)
	if err != nil {
		return false
	}
	if h.Host == "" {
		return false
	}
	return true
}
