package dcrypto

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func ParsePemOpenSSHKeyToSigner(keyPem string) (sshSigner ssh.Signer) {
	sshSigner, err := ssh.ParsePrivateKey([]byte(keyPem))
	handleError(err)
	return
}

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
