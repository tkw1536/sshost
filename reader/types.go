package reader

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Default reads a value from reader, or returns dflt if the value is empty
func Default(reader Reader, alias, key, dflt string) (string, error) {
	value, err := reader.Get(alias, key)
	if err != nil {
		return value, err
	}
	if value == "" {
		return dflt, nil
	}

	return value, nil
}

// StringSlice reads a comma-seperated string slice from reader
func StringSlice(reader Reader, alias, key string, dflt []string) ([]string, error) {
	value, err := Default(reader, alias, key, "")
	if err != nil {
		return nil, err
	}
	if value == "" {
		return dflt, nil
	}
	return strings.Split(value, ","), nil
}

// Uint reads an unsigned integer from reader
func Uint(reader Reader, alias, key string, base, bitSize int, dflt uint64) (uint64, error) {
	value, err := Default(reader, alias, key, "")
	if err != nil {
		return 0, err
	}
	if value == "" {
		return dflt, nil
	}

	return strconv.ParseUint(value, base, bitSize)
}

// Int reads an integer value from reader
func Int(reader Reader, alias, key string, base, bitSize int, dflt int64) (int64, error) {
	value, err := Default(reader, alias, key, "")
	if err != nil {
		return 0, err
	}
	if value == "" {
		return dflt, nil
	}

	return strconv.ParseInt(value, base, bitSize)
}

// Seconds reads a time.Duration value in seconds from reader
func Seconds(reader Reader, alias, key string, dflt time.Duration) (time.Duration, error) {
	value, err := Int(reader, alias, key, 10, 64, int64(dflt.Seconds()))
	if err != nil {
		return 0, err
	}
	return time.Duration(value) * time.Second, nil
}

// YesNo reads a "yes" or "no" from reader
// Returns ErrNotABoolean when reading fails.
func YesNo(Config Reader, alias, key string, dflt bool) (bool, error) {
	value, err := Default(Config, alias, key, "")
	if err != nil {
		return false, err
	}
	if value == "" {
		return dflt, nil
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, ErrNotABoolean
	}
}

// ErrNotABoolean is returned when a value is not a boolean
var ErrNotABoolean = errors.New("received non-boolean value")
