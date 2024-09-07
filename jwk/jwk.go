package jwk

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

type ECJwk struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	D   string `json:"d"`
}

var loadedPrivateKey *ecdsa.PrivateKey

func loadKey(location string) error {
	pkey := new(ecdsa.PrivateKey)
	bytes, err := os.ReadFile(location)
	if err != nil {
		log.Println("Could not load private key, generating in-memory key")
		pkey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		bytes, err = x509.MarshalECPrivateKey(pkey)
		if err != nil {
			return err
		}
		bytes = pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: bytes,
		})
	}

	block, _ := pem.Decode(bytes)
	pkey, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	loadedPrivateKey = pkey
	log.Println("Loaded private key")
	return nil
}

/*
GetSigningKey returns the private key used to sign tokens

If the key has not been loaded, it will load it from the default location
`.keys/tokenSigner`
*/
func GetSigningKey() (*ecdsa.PrivateKey, error) {
	if loadedPrivateKey == nil {
		err := loadKey("./.keys/tokenSigner")
		if err != nil {
			return nil, err
		}
	}
	return loadedPrivateKey, nil
}
