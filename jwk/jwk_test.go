package jwk_test

import (
	"authmantle-sso/jwk"
	"github.com/golang-jwt/jwt/v5"
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

func TestTokenVerification(t *testing.T) {
	var jwkKey jwk.ECJwk
	var token string
	{
		privateKey, err := jwk.GetSigningKey()
		if err != nil {
			t.Fatal(err)
		}

		idToken := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
			"sub": 1234567890,
			"iss": "testing.com",
			"aud": "https://sso.demorith.com",
			"iat": 1516239022,
		})
		token, err = idToken.SignedString(privateKey)
		if err != nil {
			t.Fatal("Failed to encode JWKs:", err)
		}

		jwkKey = jwk.GetEcJWK(privateKey)
	}
	publicKey, err := jwk.GetKeyFromJWK(jwkKey)
	if err != nil {
		t.Fatal("Failed to decode JWKs:", err)
	}
	parse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		t.Fatal("Failed to parse JWT:", err)
	}
	if parse.Valid {
		t.Log(parse.Header)
		t.Log(parse.Claims)
		t.Log(parse.Raw)
		t.Log("Token is valid")
	} else {
		t.Error("Token is invalid")
	}
}
