package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"testing"
)

func TestSerDE(t *testing.T) {
	key := ED205519KeyGen()
	bits := MarshalED205519Key(key)
	if len(bits) < 31 {
		t.Errorf("not enpugh bits %d", len(bits))
	}
	fmt.Println(string(bits))

	signer, err := ssh.ParsePrivateKey(bits)
	handleError(err)
    fmt.Println(signer.PublicKey().Type())

	signerFromString := ParsePemOpenSSHKey(string(bits))
	fmt.Println(signerFromString.PublicKey().Type())

}


func TestDeploy(t *testing.T) {
	key := ED205519KeyGen()
	bits := MarshalED205519Key(key)
	if len(bits) < 31 {
		t.Errorf("not enpugh bits %d", len(bits))
	}
	// pem key
	fmt.Println(string(bits))
	signer, err := ssh.ParsePrivateKey(bits)
	handleError(err)
	fp := ssh.FingerprintSHA256(signer.PublicKey())
	// finger print
	fmt.Println(fp)
}




const priKeyPem = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtz
c2gtZWQyNTUxOQAAACD59wgphF8jW/K+0j8kPWvDKPWOLoaUPpdeP3k+PmJoOQAA
AIOaywRCmssEQgAAAAtzc2gtZWQyNTUxOQAAACD59wgphF8jW/K+0j8kPWvDKPWO
LoaUPpdeP3k+PmJoOQAAAEBpnDxun42XGx95FI4UM1Yd42Rp0BaBLOUTDgy7jTEG
FPn3CCmEXyNb8r7SPyQ9a8Mo9Y4uhpQ+l14/eT4+Ymg5AAAAAA==
-----END OPENSSH PRIVATE KEY-----`
