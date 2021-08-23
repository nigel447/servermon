package server

import (

	"golang.org/x/crypto/ssh"
	"net"
    
)

type (
	Server struct {
		Addr           string
		ClienPublicKey ssh.PublicKey
		config *ssh.ServerConfig
	}
)

var (
	lc net.ListenConfig

)


 