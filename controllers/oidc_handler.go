package controllers

import (
	"authmantle-sso/data"
	"authmantle-sso/jwk"
	"authmantle-sso/middleware"
	"authmantle-sso/utils"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

type EndpointHelper struct {
	Method      string
	Endpoint    string
	FunctionPTR func(w http.ResponseWriter, r *http.Request)
}
type AuthRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}
type Page struct {
	PageMeta   MetaData
	StatusCode int
	Error      string
}
type MetaData struct {
	PageTitle string
}

type Controller struct {
	Db       *data.DatabaseHandler
	Renderer *utils.Renderer
}

func (c *Controller) HandleWellKnown(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	wk := new(data.WellKnownResponse) // temporary retardation.. TODO remove this shit
	err := json.NewEncoder(w).Encode(wk)
	if err != nil {
		slog.ErrorContext(ctx, "Error while encoding JWKs", "error", err)
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
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

func (c *Controller) HandleJWKs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jwkList := make([]jwk.ECJwk, 1) // TODO remove and actually parse some keys
	defer func() {
		jwkList = nil // power to the ppl bby
	}()
	privateKey, err := jwk.GetSigningKey()
	if err != nil {
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
		return
	}
	jwkList[0] = jwk.ECJwk{
		Kty: "EC",
		Crv: "P-256",
		X:   fmt.Sprintf("%x", privateKey.X),
		Y:   fmt.Sprintf("%x", privateKey.Y),
		D:   fmt.Sprintf("%x", privateKey.D),
	}
	j := data.JWKResponse[jwk.ECJwk]{Keys: &jwkList}
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

func (c *Controller) HandleNewToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	connection, err := c.Db.Acquire(ctx)
	defer connection.Release()
	if err != nil {
		logger.Error("Failed to acquire connection", "error", err)
		http.Error(w, "Failed to acquire connection", http.StatusInternalServerError)
		return
	}

	req := new(AuthRequest)
	defer func() {
		req = nil
	}()
	err = req.ParseContent(r.Header.Get("Content-Type"), r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to parse request", http.StatusInternalServerError)
		return
	}
	authCode := new(data.AuthCodeRequest)
	err = authCode.GetAuthCodeRequest(ctx, connection, req.Code)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get auth code", http.StatusInternalServerError)
		return
	}
	log.Println(req)
	// TODO validate grant_type, scopes, code and redirect_uri(again)

	res := &data.AuthResponse{
		Scope:     "openid profile email",
		ExpiresIn: 86400,
		TokenType: "Bearer",
	}
	privateKey, err := jwk.GetSigningKey()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
		return
	}

	err = authCode.ConsumeAuthCodeRequest(ctx, connection)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to consume auth code", http.StatusInternalServerError)
		return
	}

	idToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"foo": "idToken",
	})
	if token, err := idToken.SignedString(privateKey); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
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
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
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
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
		return
	}
}

func (c *Controller) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	c.Renderer.Render(w, "login.html", Page{PageMeta: MetaData{PageTitle: "Login"}})
}

func (c *Controller) HandleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	connection, err := c.Db.Acquire(ctx)
	defer connection.Release()
	if err != nil {
		logger.Error("Failed to acquire connection", "error", err)
		http.Error(w, "Failed to acquire connection", http.StatusInternalServerError)
		return
	}

	user := new(data.User)
	err = user.GetUser(ctx, *logger, connection, r.FormValue("username"))
	if err != nil {
		logger.Warn("User not found", "username", r.FormValue("username"), "error", err)
		err = c.Renderer.Render(w, "login.html", Page{
			PageMeta: MetaData{PageTitle: "Login"},
			Error:    "Invalid Password or Username",
		})
		if err != nil {
			logger.Error("Failed to render login page", "error", err)
		}
		return
	}
	if user.Password != r.FormValue("password") {
		logger.WarnContext(ctx, "User's credentials did not match!", "username", r.FormValue("username"))
		err = c.Renderer.Render(w, "login.html", Page{
			PageMeta: MetaData{PageTitle: "Login"},
			Error:    "Invalid Password or Username",
		})
		if err != nil {
			logger.ErrorContext(ctx, "Failed to render login page", "error", err)
		}
		return
	}
	redir := r.URL.Query().Get("redirect_uri")
	valid, err := data.CheckRedirectURI(ctx, connection, redir)
	if redir == "" || err != nil || !valid {
		logger.ErrorContext(ctx, "Invalid redirect_uri", "redirect_uri", redir)
		err = c.Renderer.Render(w, "login.html", Page{
			PageMeta: MetaData{PageTitle: "Login"},
			Error:    "Invalid redirect_uri",
		})
		if err != nil {
			logger.ErrorContext(ctx, "Failed to render login page", "error", err)
		}
		return
	}
	authReq := new(data.AuthCodeRequest)
	err = authReq.CreateAuthCodeRequest(ctx, connection, user.ID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to create auth code", "error", err)
		err = c.Renderer.Render(w, "login.html", Page{
			PageMeta: MetaData{PageTitle: "Login"},
			Error:    "Auth code error, please try again later",
		})
		if err != nil {
			logger.ErrorContext(ctx, "Failed to render login page", "error", err)
		}
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%s?code=%s", redir, authReq.AuthCode), http.StatusSeeOther) // hehe, stupid shit going down right here ;)
}

func (c *Controller) GetLandingPage(w http.ResponseWriter, r *http.Request) {
	if s := r.URL.Path; s != "/" { // make sure that the shit does not effect other pages.
		http.Redirect(w, r, "/error/404", http.StatusSeeOther)
		return
	}
	c.Renderer.Render(w, "index.html", Page{PageMeta: MetaData{PageTitle: "Login"}})
}
func (c *Controller) GetRegisterPage(w http.ResponseWriter, r *http.Request) {
	c.Renderer.Render(w, "register.html", Page{PageMeta: MetaData{PageTitle: "Login"}})
}
func (c *Controller) GetAdminPage(w http.ResponseWriter, r *http.Request) {
	c.Renderer.Render(w, "admin_login.html", Page{PageMeta: MetaData{PageTitle: "Admin Login"}})
}
func (c *Controller) ErrorRedirect(w http.ResponseWriter, r *http.Request) {
	status := parseStatusCode(r.PathValue("status"))
	c.Renderer.Render(w, "error.html", Page{PageMeta: MetaData{PageTitle: "Error"}, StatusCode: status})
}
func parseStatusCode(pathError string) int {
	if pathError == "" {
		return http.StatusInternalServerError
	}
	if len(pathError) > 4 {
		return http.StatusInternalServerError
	}
	status, err := strconv.Atoi(pathError)
	if err != nil {
		return http.StatusInternalServerError
	}

	return status
}

func (c *Controller) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	c.Renderer.Render(w, "user_settings.html", Page{PageMeta: MetaData{PageTitle: "User Settings"}})
}
func (c *Controller) GetAdminDashboardPage(w http.ResponseWriter, r *http.Request) {
	c.Renderer.Render(w, "admin_panel.html", Page{PageMeta: MetaData{PageTitle: "Admin Dashboard"}})
}
