package dcrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

 
func GenECDSAKey() (*ecdsa.PrivateKey, ecdsa.PublicKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, ecdsa.PublicKey{}, err
	}
	return priv, priv.PublicKey, nil
}
 