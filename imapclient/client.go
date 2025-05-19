package imapclient

import (
	"crypto/tls"

	"github.com/emersion/go-imap/client"
)

func Connect(email, password, host string) (*client.Client, error) {
	c, err := client.DialTLS(host+":993", &tls.Config{})
	if err != nil {
		return nil, err
	}

	if err := c.Login(email, password); err != nil {
		return nil, err
	}

	return c, nil
}
