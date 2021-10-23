package sshost

import (
	"github.com/tkw1536/sshost/pkg/expand"
)

func (profile Profile) Expander() expand.Expander {
	return expand.Expander{
		Getenv: profile.env.getenv,
	}
}

var identityFileFlags = expand.Flags{
	Environment: true,
	Tilde:       true,
	Tokens:      "%CdhikLlnpru",
}

func (profile Profile) identityFileExpander() expand.Expander {
	return profile.Expander()
}

// IdentityFile returns the IdentityFile being used by this profile
func (profile Profile) IdentityFile() []string {
	result := make([]string, 0, len(profile.config.IdentityFile))
	ex := profile.identityFileExpander()
	for _, id := range profile.config.IdentityFile {
		name, err := ex.Expand(id, identityFileFlags)
		if err != nil {
			continue
		}
		result = append(result, name)
	}
	return result
}

var identityAgentFlags = expand.Flags{
	Environment: true,
	Tilde:       true,
	Tokens:      "%CdhikLlnpru",
}

func (profile Profile) identityAgentExpander() expand.Expander {
	return profile.Expander()
}

// IdentityAgent returns the identity agent to connect to
func (profile Profile) IdentityAgent() string {
	agent := profile.config.IdentityAgent

	if agent == "none" || agent == "" {
		return ""
	}
	if agent == "SSH_AUTH_SOCK" {
		return profile.env.getenv("SSH_AUTH_SOCK")
	}
	if agent[0] == '$' {
		return profile.env.getenv(agent[1:])
	}

	ex := profile.identityAgentExpander()
	agent, _ = ex.Expand(agent, identityAgentFlags)
	return agent
}
