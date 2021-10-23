package expand

import (
	"fmt"
	"strings"
)

func ExampleExpander_Expand() {
	// create an expansion
	ex := Expander{
		Getenv: func(name string) string {
			switch strings.ToUpper(name) {
			case "SOMETHING":
				return "something-secret"
			case "HOME":
				return "/path/to/home"
			default:
				return ""
			}
		},
	}

	// expand an actual string
	expanded, err := ex.Expand("${HOME}/%%/${SOMETHING}", Flags{Tokens: AllTokens, Environment: true})
	if err != nil {
		panic(err)
	}
	fmt.Println(expanded)

	untouched, err := ex.Expand("${HOME}/%%/${SOMETHING}", Flags{})
	if err != nil {
		panic(err)
	}
	fmt.Println(untouched)

	// Output:
	// /path/to/home/%/something-secret
	// ${HOME}/%%/${SOMETHING}

}
