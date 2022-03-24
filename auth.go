package sshost

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

// AuthEnv is an environment to run authentication methods in
type AuthEnv struct {
	Stdin  io.Reader
	Stdout io.Writer

	PasswordPrompt string
}

// in returns the input used for this environment
func (m AuthEnv) in() io.Reader {
	if m.Stdin == nil {
		return os.Stdin
	}
	return m.Stdin
}

// inFD returns the file descriptor for the input
func (m AuthEnv) inFD() int {
	switch inT := m.in().(type) {
	case *os.File:
		return int(inT.Fd())
	default:
		return syscall.Stdin
	}
}

// out returns the output used for this environment
func (m AuthEnv) out() io.Writer {
	if m.Stdout == nil {
		return os.Stdout
	}
	return m.Stdout
}

// write writes message to the output
func (m AuthEnv) print(message string, newline bool) {
	if newline {
		fmt.Fprintln(m.out(), message)
	} else {
		fmt.Fprint(m.out(), message)
	}
}

// readOpen openly reads text from standard input
func (m AuthEnv) readOpen() (string, error) {
	reader := bufio.NewReader(m.in())
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// readClosed reads text from standard input, hiding the output.
func (m AuthEnv) readClosed() (string, error) {
	bytes, err := term.ReadPassword(m.inFD())
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(bytes)), nil
}

const DefaultPasswordPrompt = "Password: "

type AuthMethod string

const (
	GssapiWithMic       AuthMethod = "gssapi-with-mic" // unsupported
	HostBased           AuthMethod = "hostbased"       // unsupported
	PublicKey           AuthMethod = "publickey"
	KeyboardInteractive AuthMethod = "keyboard-interactive"
	Password            AuthMethod = "password"
)

// Methods gets all ssh methods specified in the comma-seperated string methods.
// When a method does not exist, or is not supported, skips over it.
func (m AuthEnv) Methods(methods string, profile Profile) []ssh.AuthMethod {
	names := strings.Split(methods, ",")
	auths := make([]ssh.AuthMethod, 0)
	for _, name := range names {
		auth := m.Method(AuthMethod(name), profile)
		auths = append(auths, auth...)
	}
	return auths
}

// Method gets the specified authentication method.
// When the method does not exist, or is not supported, returns nil.
func (m AuthEnv) Method(method AuthMethod, profile Profile) []ssh.AuthMethod {
	switch method {
	case PublicKey:
		return m.mPublicKey(profile)
	case KeyboardInteractive:
		return m.mKeyboardInteractive(profile)
	case Password:
		return m.mPassword(profile)
	}
	return nil
}

// mPassword returns the password authentication method.
func (m AuthEnv) mPassword(profile Profile) []ssh.AuthMethod {
	if !profile.config.PasswordAuthentication {
		return nil
	}
	return []ssh.AuthMethod{ssh.RetryableAuthMethod(
		ssh.PasswordCallback(func() (secret string, err error) {
			return m.passwordCallback()
		}),
		profile.config.NumberOfPasswordPrompts,
	)}
}

// passwordCallback is invoked to retrieve the password.
func (m AuthEnv) passwordCallback() (secret string, err error) {
	prompt := m.PasswordPrompt
	if prompt == "" {
		prompt = DefaultPasswordPrompt
	}

	m.print(prompt, false)
	return m.readClosed()
}

// mKeyboardInteractive returns the keyboard-interactive authentication method.
func (m AuthEnv) mKeyboardInteractive(profile Profile) []ssh.AuthMethod {
	if !profile.config.KbdInteractiveAuthentication {
		return nil
	}

	return []ssh.AuthMethod{ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
		return m.keyboardInteractiveCallback(user, instruction, questions, echos)
	})}
}

// keyboardInteractiveCallback is invoked to make use of the keyboard-interactive method
func (m AuthEnv) keyboardInteractiveCallback(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	m.print(user, true)
	m.print(instruction, true)

	answers = make([]string, len(questions))
	for i, q := range questions {
		m.print(q, false)
		if echos[i] {
			answers[i], err = m.readOpen()
		} else {
			answers[i], err = m.readClosed()
		}
		if err != nil {
			return nil, err
		}
	}
	return answers, nil
}

// mPublicKey returns the public-key authentication method
func (m AuthEnv) mPublicKey(profile Profile) []ssh.AuthMethod {
	methods := make([]ssh.AuthMethod, 0)
	if agent := m.publicKeyAgent(profile); agent != nil {
		methods = append(methods, agent)
	}
	IdentityFile := profile.IdentityFile()
	for _, file := range IdentityFile {
		if pk := m.identityFile(file); pk != nil {
			methods = append(methods, pk)
		}
	}
	return methods
}

func (m AuthEnv) identityFile(path string) ssh.AuthMethod {
	// read the bytes, it's fine if we can't read the file.
	pkBytes, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	// decode the bytes as a public key, but error out if they can't be read!
	return ssh.PublicKeysCallback(func() (signers []ssh.Signer, err error) {
		signer, err := ssh.ParsePrivateKey(pkBytes)
		if err != nil {
			return nil, err
		}
		return []ssh.Signer{signer}, nil
	})
}

func (m AuthEnv) publicKeyAgent(profile Profile) ssh.AuthMethod {
	IdentityAgent := profile.IdentityAgent()
	if profile.config.IdentitiesOnly || IdentityAgent == "" {
		return nil
	}

	// read the identity agent!
	return ssh.PublicKeysCallback(func() (signers []ssh.Signer, err error) {
		agentc, err := net.Dial("unix", IdentityAgent)
		if err != nil {
			return nil, err
		}
		client := agent.NewClient(agentc)
		return client.Signers()
	})
}
