package sshclient

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

type (
	Client struct {
		Addr           string
		ClienPublicKey ssh.PublicKey
		config *ssh.ClientConfig
	}
)


func FingerprintKey(k ssh.PublicKey) string {
	bytes := md5.Sum(k.Marshal())
	strbytes := make([]string, len(bytes))
	for i, b := range bytes {
		strbytes[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(strbytes, ":")
}

func s(f string, args ...interface{}) []byte {
	return []byte(fmt.Sprintf(f, args...))
}
