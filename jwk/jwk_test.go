package jwk_test

import (
	"authmantle-sso/jwk"
	"github.com/google/uuid"
	"log"
	"testing"
)

func BenchmarkGetSigningKey(b *testing.B) {
	var err error
	uuid.EnableRandPool()
	for i := 0; i < b.N; i++ {
		_, err = jwk.GetSigningKey()
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Printf("BENCHMARK: Loaded private key %d times", b.N)
}
