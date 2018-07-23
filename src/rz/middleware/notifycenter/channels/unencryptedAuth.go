package channels

import (
	"net/smtp"
	"errors"
)

type unencryptedAuth struct {
	smtp.Auth
}

func (unencryptedAuth *unencryptedAuth) Start(serverInfo *smtp.ServerInfo) (string, []byte, error) {
	(*serverInfo).TLS = true
	return unencryptedAuth.Auth.Start(serverInfo)
}

func (*unencryptedAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// We've already sent everything.
		return nil, errors.New("unexpected server challenge")
	}

	return nil, nil
}