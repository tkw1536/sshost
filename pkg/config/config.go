// Package config contains Config and related datastructures.
//
// All types in this configuration are considered stateless.
package config

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/tkw1536/stringreader"
)

// Config represents a configuration for a single host
type Config struct {
	AddressFamily       AddressFamily `config:"AddressFamily" type:"string"`
	Ciphers             []string      `config:"Ciphers" type:"stringslice"`
	Compression         bool          `config:"Compression" type:"yesno"`
	ConnectionAttempts  uint64        `config:"ConnectionAttempts" type:"int"`
	ConnectTimeout      time.Duration `config:"ConnectTimeout" type:"seconds"`
	HostKeyAlgorithms   []string      `config:"HostKeyAlgorithms" type:"stringslice"`
	Hostname            string        `config:"Hostname" type:"string"`
	IdentityAgent       string        `config:"IdentityAgent" type:"string"`
	KexAlgorithms       []string      `config:"KexAlgorithms" type:"stringslice"`
	MACs                []string      `config:"MACs" type:"stringslice"`
	ProxyJump           []string      `config:"ProxyJump" type:"stringslices"` // TODO: multi-slice
	Port                uint16        `config:"Port" type:"uint"`
	RekeyLimit          string        `config:"RekeyLimit" type:"string"` // TODO: Proper datatype
	ServerAliveCountMax uint64        `config:"ServerAliveCountMax" type:"uint"`
	ServerAliveInterval time.Duration `config:"ServerAliveInterval" type:"seconds"`
	Username            string        `config:"Username" type:"string"`
}

// Defaults contains defaults for generating an environment
type Defaults struct {
	Username string
}

func (dflts Defaults) Data() (data stringreader.ParsingData) {
	data.SetLocal("Hostname", "default", "")

	data.SetLocal("Port", "default", 22)
	data.SetLocal("Port", "base", 10)
	data.SetLocal("Port", "bits", 16)

	data.SetLocal("Username", "default", dflts.Username)

	data.SetLocal("IdentityAgent", "default", "$SSH_AUTH_SOCK")

	data.SetLocal("AddressFamily", "default", string(DefaultAddressFamily))

	data.SetLocal("HostKeyAlgorithms", "default", nil)

	data.SetLocal("Ciphers", "default", nil)

	data.SetLocal("Compression", "default", false)

	data.SetLocal("ConnectionAttempts", "default", 1)
	data.SetLocal("ConnectionAttempts", "base", 10)
	data.SetLocal("ConnectionAttempts", "bits", 64)

	data.SetLocal("ConnectTimeout", "default", time.Second)

	data.SetLocal("KexAlgorithms", "default", nil)

	data.SetLocal("MACs", "default", nil)

	data.SetLocal("ServerAliveInterval", "default", 0)

	data.SetLocal("ServerAliveCountMax", "default", 3)
	data.SetLocal("ServerAliveCountMax", "base", 10)
	data.SetLocal("ServerAliveCountMax", "bits", 64)

	data.SetLocal("RekeyLimit", "default", "default none")

	data.SetLocal("ProxyJump", "default", nil)
	data.SetLocal("ProxyJump", "skip", "none")

	return
}

// NewConfig reads a configuration from the provided source.
func NewConfig(source stringreader.Source, host Host, dflts Defaults) (cfg Config, err error) {
	if err = checkUnsupportedConfig(source); err != nil {
		return
	}
	if err = configMarshal.UnmarshalContext(&cfg, source, dflts.Data()); err != nil {
		return
	}
	if err = host.UpdateConfig(&cfg); err != nil {
		return
	}
	return
}

// UpdateConfig updates cfg with configuration values from this host.
// Always returns nil.
func (h Host) UpdateConfig(cfg *Config) error {
	if cfg.Hostname == "" {
		cfg.Hostname = h.Host
	}

	if h.User != "" {
		cfg.Username = h.User
	}
	if h.Port != 0 {
		cfg.Port = h.Port
	}
	return nil
}

var configMarshal stringreader.Marshal

func init() {
	configMarshal.NameTag = "config"
	configMarshal.StrictNameTag = true
	configMarshal.ParserTag = "type"
	configMarshal.DefaultParser = ""

	configMarshal.RegisterSingleParser("string", func(value string, ok bool, ctx stringreader.ParsingContext) (interface{}, error) {
		if !ok || value == "" {
			return ctx.Get("default"), nil
		}
		return value, nil
	})

	configMarshal.RegisterSingleParser("stringslice", func(value string, ok bool, ctx stringreader.ParsingContext) (interface{}, error) {
		if !ok {
			return ctx.Get("default"), nil
		}
		if value == "" {
			return nil, nil
		}
		return strings.Split(value, ","), nil
	})
	configMarshal.RegisterMultiParser("stringslices", func(values []string, ok bool, ctx stringreader.ParsingContext) (interface{}, error) {
		if !ok {
			return ctx.Get("default"), nil
		}
		var results []string
		for _, value := range values {
			if value == "" {
				continue
			}
			results = append(results, strings.Split(value, ",")...)
		}

		// remove skipped values!
		skip := ctx.Get("skip")
		if sskip, ok := skip.(string); ok {
			nresults := results[:0]
			for _, value := range results {
				if value == sskip {
					continue
				}
				nresults = append(nresults, value)
			}
			results = nresults
		}

		return results, nil
	})

	configMarshal.RegisterSingleParser("int", func(value string, ok bool, ctx stringreader.ParsingContext) (interface{}, error) {
		if !ok || value == "" {
			return ctx.Get("default"), nil
		}
		return strconv.ParseInt(value, ctx.Get("base").(int), ctx.Get("bits").(int))
	})

	configMarshal.RegisterSingleParser("uint", func(value string, ok bool, ctx stringreader.ParsingContext) (interface{}, error) {
		if !ok || value == "" {
			return ctx.Get("default"), nil
		}
		return strconv.ParseUint(value, ctx.Get("base").(int), ctx.Get("bits").(int))
	})

	configMarshal.RegisterSingleParser("seconds", func(value string, ok bool, ctx stringreader.ParsingContext) (interface{}, error) {
		if !ok || value == "" {
			return ctx.Get("default"), nil
		}
		s, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return 0, err
		}
		return time.Duration(s) * time.Second, nil
	})

	configMarshal.RegisterSingleParser("yesno", func(value string, ok bool, ctx stringreader.ParsingContext) (interface{}, error) {
		if !ok || value == "" {
			return ctx.Get("default"), nil
		}

		switch strings.ToLower(strings.TrimSpace(value)) {
		case "yes":
			return true, nil
		case "no":
			return false, nil
		default:
			return false, ErrNotABoolean
		}
	})
}

// ErrNotABoolean is returned when a value is not a boolean
var ErrNotABoolean = errors.New("received non-boolean value")
