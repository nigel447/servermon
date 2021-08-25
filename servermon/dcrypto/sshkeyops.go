package dcrypto

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func GenOpenSSHSigner() (sshSigner ssh.Signer) {
	pKey, err := ssh.ParsePrivateKey([]byte(PriKeyPem))
	handleError(err)
	sshSigner = pKey
	return
}

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
