package client

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type (
	Client struct {
		Addr           string
		ClienPublicKey ssh.PublicKey
		config         *ssh.ClientConfig
	}
)

func s(f string, args ...interface{}) []byte {
	return []byte(fmt.Sprintf(f, args...))
}
