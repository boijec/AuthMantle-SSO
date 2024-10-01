package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// TODO add ability to generate key with passphrase
func main() {
	if len(os.Args) == 0 {
		panic("No arguments provided")
	}
	fileName := os.Args[1]
	private, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	privateBytes, _ := x509.MarshalECPrivateKey(private)
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	pem.Encode(f, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateBytes,
	})
}
