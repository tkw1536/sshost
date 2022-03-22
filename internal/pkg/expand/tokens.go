package expand

import (
	"fmt"
	"strings"
)

// TokenList is a list of tokens to expand
type TokenList string

func (t TokenList) Any() bool {
	return len(string(t)) > 0
}

// Allowed checks if a specific token can be expanded by this token list
func (t TokenList) Allowed(r rune) bool {
	return strings.ContainsRune(string(t), r)
}

// AllTokens is a list of all supported tokens
const AllTokens TokenList = "%"

func (ex Expander) ExpandToken(r rune) (string, error) {
	switch r {
	case '%':
		return "%", nil
	default:
		return "", fmt.Errorf("encounted unknown/unimplemented '%%' token: %q", r)
	}
}
