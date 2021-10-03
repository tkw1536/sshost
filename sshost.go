package sshost

import "golang.org/x/crypto/ssh"

// NewSession creates a new session in the provided environment
func NewSession(ctx *Context, alias string) (*ssh.Session, *ClosableStack, error) {
	profile, err := ctx.NewProfile(alias)
	if err != nil {
		return nil, nil, err
	}

	conn, closers, err := profile.Dial()
	if err != nil {
		return nil, nil, err
	}

	session, err := profile.Connect(conn)
	if err != nil {
		defer closers.Close()
		return nil, nil, err
	}

	closers.Push(session)
	return session, closers, err
}
