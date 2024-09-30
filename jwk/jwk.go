package jwk

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"sync"
)

type ECJwk struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Use string `json:"use"`
}

type LoadedKey struct {
	lock sync.Mutex
	key  *ecdsa.PrivateKey
}

var loadedPrivateKey *LoadedKey

func loadKey(location string) error {
	loadedPrivateKey = new(LoadedKey)
	loadedPrivateKey.lock = sync.Mutex{}
	loadedPrivateKey.lock.Lock()
	defer func() {
		loadedPrivateKey.lock.Unlock()
	}()

	pkey := new(ecdsa.PrivateKey)
	bytes, err := os.ReadFile(location)
	if err != nil {
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
	loadedPrivateKey.key = pkey
	log.Println("Loaded private key")
	return nil
}

/*
GetSigningKey returns the private key used to sign tokens

If the key has not been loaded, it will load it from the default location
`.keys/tokenSigner` or create an in-memory key if the file does not exist
*/
func GetSigningKey() (*ecdsa.PrivateKey, error) {
	if loadedPrivateKey == nil {
		err := loadKey("./.keys/tokenSigner")
		if err != nil {
			return nil, err
		}
	}
	loadedPrivateKey.lock.Lock()
	defer loadedPrivateKey.lock.Unlock()
	return loadedPrivateKey.key, nil
}
