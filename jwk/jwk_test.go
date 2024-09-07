package jwk_test

import (
	"authmantle-sso/jwk"
	"log"
	"testing"
)

func BenchmarkGetSigningKey(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = jwk.GetSigningKey()
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Printf("BENCHMARK: Loaded private key %d times", b.N)
}
