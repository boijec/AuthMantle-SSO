package oidc

import (
	"encoding/json"
	"fmt"
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
type JWKResponse struct {
	Keys *[]JWK `json:"keys"`
}
type JWK struct {
	Kid string    `json:"kid"`
	Kty string    `json:"kty"`
	Alg string    `json:"alg"`
	Use string    `json:"use"`
	N   string    `json:"n"`
	E   string    `json:"e"`
	X5t string    `json:"x5t"`
	X5c [1]string `json:"x5c"`
}

type AuthRequest struct {
	GrantType   string
	Code        string
	RedirectUri string
}
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
type EndpointHelper struct {
	Method      string
	Endpoint    string
	FunctionPTR func(w http.ResponseWriter, r *http.Request)
}

// thank FUCK for globals - but this needs some more thought
var ConfiguredRoutes = map[string]EndpointHelper{
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
		"response_modes_supported": ["query", "fragment"],
		"grant_types_supported": ["authorization_code"],
		"subject_types_supported": ["public"],
		"id_token_signing_alg_values_supported": ["RS256"],
		"token_endpoint_auth_methods_supported": ["client_secret_basic"],
		"claims_supported": ["sub", "iss", "email", "profile"],
		"code_challenge_methods_supported": ["plain", "S256"]
	*/
}

func HandleJWKs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jwkList := make([]JWK, 2) // TODO remove and actually parse some keys
	jwk := JWKResponse{Keys: &jwkList}
	defer func() {
		jwkList = nil // power to the ppl bby
	}()
	err := json.NewEncoder(w).Encode(jwk)
	if err != nil {
		http.Error(w, "Failed to encode jwkList", http.StatusInternalServerError)
		return
	}
}

func HandleNewToken(w http.ResponseWriter, r *http.Request) {
	// for good results load from db, for best results hard code the shit out of it? TODO check if this can be dynamically checked (support more than form data)
	req := &AuthRequest{
		GrantType:   r.FormValue("grant_type"),
		Code:        r.FormValue("code"),
		RedirectUri: r.FormValue("redirect_uri"),
	}
	log.Println(req)

	res := &AuthResponse{
		Scope:     "openid profile email",
		ExpiresIn: 86400,
		TokenType: "Bearer",
	}
	res.AccessToken = "accessToken"
	res.IdToken = "idToken"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200 OK")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Failed to encode jwks", http.StatusInternalServerError)
		return
	}
	// should respond with
	// 	"{"access_token": "accessToken","id_token":"idToken","scope":"openid profile email","expires_in":86400,"token_type":"Bearer"}
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	redir := r.URL.Query().Get("redirect_uri")
	http.Redirect(w, r, fmt.Sprintf("%s?code=%s", redir, "dudde1234"), http.StatusSeeOther) // hehe, stupid shit going down right here ;)
}
