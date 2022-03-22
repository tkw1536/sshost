// Package expand implements expansion of environment variables
package expand

import (
	"errors"
	"fmt"
	"strings"
)

// Expander expands tokens, tildes and environment variables.
// This struct holds context required for expansion.
type Expander struct {
	Getenv func(string) string
}

// Flags determines which expands an Expander should perform.
type Flags struct {
	// List of tokens to expand
	Tokens TokenList

	// When true, expand environment variables
	Environment bool

	// When true, expand ~ into ${HOME} variable
	Tilde bool
}

var ErrTrailingPercent = errors.New("Expander.Expand: Encountered trailing '%' sign")
var ErrUnclosedVariable = errors.New("Expander.Expand: Unclosed variable")

// Expand expands value as defined in flags.
func (ex Expander) Expand(value string, flags Flags) (string, error) {
	var result strings.Builder

	var modePercent bool
	var modeEnv bool
	var modeEnvReadBrace bool
	var modeEnvGobble strings.Builder

	rValue := []rune(value)
	for i, r := range rValue {
		// percent mode: replace the current '%' character
		// with the replacement!
		if modePercent {
			value, err := ex.ExpandToken(r)
			if err != nil {
				return "", err
			}
			result.WriteString(value)
			modePercent = false
			continue
		}

		// environment mode: first check that we found a '{',
		// then gobble until trailing '}'
		if modeEnv {
			switch {
			case !modeEnvReadBrace:
				// this check should never trigger, so it's removed
				//if r != '{' {
				//	panic("Encountered ${, but did not read {")
				//}
				modeEnvReadBrace = true
			// we did not find a closing brace
			// so write it to the string
			case r != '}':
				modeEnvGobble.WriteRune(r)

			// found the closing brace
			default:
				// get the gobbled environment name
				name := modeEnvGobble.String()
				modeEnvGobble.Reset()

				// do the replacement, then go back into normal mode!
				value, err := ex.ExpandVariable(name)
				if err != nil {
					return "", err
				}
				result.WriteString(value)
				modeEnv = false
				modeEnvReadBrace = false
			}
			continue
		}

		// the cases below correspond to regular mode.

		// intercept the '%' sign for known tokens!
		if flags.Tokens.Any() && r == '%' {
			if len(rValue) > i+1 && flags.Tokens.Allowed(rValue[i+1]) {
				modePercent = true
				continue
			}
		}

		// intercept variables, we need a '${' and read until the closing '}'.
		if flags.Environment && r == '$' {
			if len(rValue) > i+1 && rValue[i+1] == '{' {
				modeEnv = true
				continue
			}
		}

		// intercept '~' when requested
		if flags.Tilde && r == '~' {
			value, err := ex.ExpandVariable("HOME")
			if err != nil {
				return "", err
			}
			result.WriteString(value)
			continue
		}

		// every other case: keep the rune!
		result.WriteRune(r)
	}

	if modePercent {
		return "", ErrTrailingPercent
	}
	if modeEnv {
		return "", ErrUnclosedVariable
	}

	return result.String(), nil
}

func (ex Expander) ExpandVariable(name string) (string, error) {
	value := ex.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("Expander.Expand: environment variable %q is unset", name)
	}
	return value, nil
}
