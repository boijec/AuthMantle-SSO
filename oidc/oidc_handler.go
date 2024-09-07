package oidc

import (
	"authmantle-sso/jwk"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
)

type WellKnownResponse struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	EndSessionEndpoint                string   `json:"end_session_endpoint"`
	JWKsUri                           string   `json:"jwks_uri"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	ResponseModesSupported            []string `json:"response_modes_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IdTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                   []string `json:"claims_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
}
type JWKResponse[T any] struct {
	Keys *[]T `json:"keys"`
}

type AuthRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}
type AuthResponse struct {
	AccessToken *string `json:"access_token"`
	IdToken     *string `json:"id_token"`
	Scope       string  `json:"scope"`
	ExpiresIn   int     `json:"expires_in"`
	TokenType   string  `json:"token_type"`
}
type EndpointHelper struct {
	Method      string
	Endpoint    string
	FunctionPTR func(w http.ResponseWriter, r *http.Request)
}

// ConfiguredRoutes global map of configured routes for OIDC discovery
var ConfiguredRoutes = map[string]*EndpointHelper{
	"jwks":  {"GET", "/.well-known/jwks.json", HandleJWKs},
	"auth":  {"POST", "/authorize", HandleAuth},
	"token": {"POST", "/oauth/token.json", HandleNewToken},
}

func HandleWellKnown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	wk := new(WellKnownResponse) // temporary retardation.. TODO remove this shit
	err := json.NewEncoder(w).Encode(wk)
	if err != nil {
		http.Error(w, "Failed to encode jwks", http.StatusInternalServerError)
		return
	}
	/*
		"issuer": "http://localhost:8080",
		"authorization_endpoint": "http://localhost:8080/v1/authorize",
		"token_endpoint": "http://localhost:8080/v1/auth/token",
		"userinfo_endpoint": "http://localhost:8080/protected/userinfo",
		"end_session_endpoint": "http://localhost:8080/v1/logout",
		"jwks_uri": "http://localhost:8080/v1/jwks.json",
		"scopes_supported": ["openid", "profile", "email"],
		"response_types_supported": ["code"],
		"grant_types_supported": ["authorization_code"],
		"subject_types_supported": ["public"],
		"id_token_signing_alg_values_supported": ["ES256"],
		"claims_supported": ["sub", "iss", "email", "profile"]
	*/
}

func HandleJWKs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jwkList := make([]jwk.ECJwk, 1) // TODO remove and actually parse some keys
	defer func() {
		jwkList = nil // power to the ppl bby
	}()
	privateKey, err := jwk.GetSigningKey()
	if err != nil {
		http.Error(w, "Failed to encode jwks", http.StatusInternalServerError)
		return
	}
	jwkList[0] = jwk.ECJwk{
		Kty: "EC",
		Crv: "P-256",
		X:   fmt.Sprintf("%x", privateKey.X),
		Y:   fmt.Sprintf("%x", privateKey.Y),
		D:   fmt.Sprintf("%x", privateKey.D),
	}
	j := JWKResponse[jwk.ECJwk]{Keys: &jwkList}
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		http.Error(w, "Failed to encode jwkList", http.StatusInternalServerError)
		return
	}
}

type ContentTypeParser interface {
	ParseContent(s string, v *http.Request) error
}

func (ar *AuthRequest) ParseContent(contentType string, req *http.Request) error {
	if req == nil {
		return fmt.Errorf("nil reference for Request")
	}
	if contentType == "" {
		return fmt.Errorf("empty Content-Type header")
	}
	switch contentType {
	case "application/x-www-form-urlencoded":
		ar.GrantType = req.FormValue("grant_type")
		ar.Code = req.FormValue("code")
		ar.RedirectUri = req.FormValue("redirect_uri")
	case "application/json":
		err := json.NewDecoder(req.Body).Decode(ar)
		if err != nil {
			return fmt.Errorf("failed to decode json body: %v", err)
		}
	default:
		return fmt.Errorf("unsupported Content-Type header")
	}
	return nil
}

func HandleNewToken(w http.ResponseWriter, r *http.Request) {
	req := new(AuthRequest)
	defer func() {
		req = nil
	}()
	err := req.ParseContent(r.Header.Get("Content-Type"), r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to parse request", http.StatusInternalServerError)
		return
	}
	log.Println(req)
	// TODO validate grant_type, scopes, code and redirect_uri(again)

	res := &AuthResponse{
		Scope:     "openid profile email",
		ExpiresIn: 86400,
		TokenType: "Bearer",
	}
	// TODO replace with fetch from globally loaded key
	// TODO cleanup
	privateKey, err := jwk.GetSigningKey()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode jwks", http.StatusInternalServerError)
		return
	}

	idToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"foo": "idToken",
	})
	if token, err := idToken.SignedString(privateKey); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode jwks", http.StatusInternalServerError)
		return
	} else {
		res.IdToken = &token
		idToken = nil
	}
	defer func() {
		res.IdToken = nil
	}()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"foo": "accessToken",
	})
	if token, err := accessToken.SignedString(privateKey); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode jwks", http.StatusInternalServerError)
		return
	} else {
		res.AccessToken = &token
		accessToken = nil
	}
	defer func() {
		res.AccessToken = nil
	}()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200 OK")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Failed to encode jwks", http.StatusInternalServerError)
		return
	}
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	// TODO check if this is a valid user or return error html
	log.Println(r.FormValue("username"), r.FormValue("password"))
	// TODO check if this is a configured redirect_uri
	redir := r.URL.Query().Get("redirect_uri")
	// TODO create a auth_code_request and store it in the db
	http.Redirect(w, r, fmt.Sprintf("%s?code=%s", redir, "dudde1234"), http.StatusSeeOther) // hehe, stupid shit going down right here ;)
}
