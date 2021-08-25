package main

// Portions of this file are derived from (https://www.scaleft.com/).
// https://github.com/ScaleFT/sshkeys/blob/master/marshal.go

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	mrand "math/rand"
)

type (
	ED25519 struct {
		Pub     []byte
		Priv    []byte
		Comment string
		Pad     []byte `ssh:"rest"`
	}

	 SshKey struct {
		Check1  uint32
		Check2  uint32
		Keytype string
		Rest    []byte `ssh:"rest"`
	}

	 Header struct {
		CipherName   string
		KdfName      string
		KdfOpts      string
		NumKeys      uint32
		PubKey       string
		PrivKeyBlock string
	}
)

const KeyTypeVersion = "openssh-key-v1"
const KeyType = "OPENSSH PRIVATE KEY"

func ED205519KeyGen() (pvk ed25519.PrivateKey) {
	_, pvk, err := ed25519.GenerateKey(rand.Reader)
	handleError(err)
	return
}

func MarshalED205519Key(key ed25519.PrivateKey) (pemBytes []byte) {
	out := Header{
		CipherName: "none",
		KdfName:    "none",
		KdfOpts:    "",
		NumKeys:    1,
		PubKey:     "",
	}

	k := ED25519{
		Pub:  key.Public().(ed25519.PublicKey),
		Priv: key,
	}
	data := ssh.Marshal(k)
	check := mrand.Uint32()
	pk1 := SshKey {
		Check1: check,
		Check2: check,
	}
	pk1.Keytype = ssh.KeyAlgoED25519
	pk1.Rest = data
	publicKey, err := ssh.NewPublicKey(key.Public())
	handleError(err)
	out.PubKey = string(publicKey.Marshal())
	out.PrivKeyBlock = string(ssh.Marshal(pk1))
	outBytes := []byte(KeyTypeVersion)
	outBytes = append(outBytes, 0)
	outBytes = append(outBytes, ssh.Marshal(out)...)

	block := &pem.Block{
		Type:  KeyType,
		Bytes: outBytes,
	}
	pemBytes, err = pem.EncodeToMemory(block), nil
	handleError(err)
	return
}

func ParsePemOpenSSHKey(keyPem string) (sshSigner ssh.Signer) {
	sshSigner, err := ssh.ParsePrivateKey([]byte(keyPem))
	handleError(err)
	return
}

func handleError(err error) {
	if err != nil {
		fmt.Println("HandleError ", err.Error())
	}
}
