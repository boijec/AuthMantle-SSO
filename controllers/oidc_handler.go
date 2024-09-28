package controllers

import (
	"authmantle-sso/data"
	"authmantle-sso/jwk"
	"authmantle-sso/middleware"
	"authmantle-sso/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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
	ClientId    string `json:"client_id"`
}

type Controller struct {
	Db       *data.DatabaseHandler
	Renderer *utils.Renderer
}

func (c *Controller) HandleWellKnown(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	realmId := ctx.Value(middleware.RealmIDContextKey).(int)
	connection, err := c.Db.Acquire(ctx)
	defer connection.Release()
	if err != nil {
		logger.Error("Failed to acquire connection", "error", err)
		http.Error(w, "Failed to acquire connection", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	rs := new(data.RealmCacheObject)
	err = rs.GetRealmSettings(ctx, connection, realmId)
	if err != nil {
		slog.ErrorContext(ctx, "Error while getting realm settings", "error", err)
		http.Error(w, "Error while getting realm settings", http.StatusInternalServerError)
		return
	}
	wk := new(data.WellKnownResponse)
	wk.ClaimsSupported = rs.Claims
	wk.GrantTypesSupported = rs.GrantTypes
	wk.ScopesSupported = rs.Scopes
	wk.SubjectTypesSupported = rs.SubjectTypes
	wk.ResponseTypesSupported = rs.ResponseTypes
	wk.IdTokenSigningAlgValuesSupported = rs.TokenSigningAlgs

	err = json.NewEncoder(w).Encode(wk)
	if err != nil {
		slog.ErrorContext(ctx, "Error while encoding JWKs", "error", err)
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
		return
	}
	// TODO implement this tomorrow, too tired to look at this for another weekend, plz god!
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
		ar.ClientId = req.FormValue("client_id")
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
		logger.ErrorContext(ctx, "Failed to acquire connection", "error", err)
		http.Error(w, "Failed to acquire connection", http.StatusInternalServerError)
		return
	}

	req := new(AuthRequest)
	err = req.ParseContent(r.Header.Get("Content-Type"), r)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to parse request", "error", err)
		http.Error(w, "Failed to parse request", http.StatusInternalServerError)
		return
	}
	authCode := new(data.AuthCodeRequest)
	err = authCode.GetAuthCodeRequest(ctx, connection, req.Code)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get auth code", "error", err)
		http.Error(w, "Failed to get auth code", http.StatusInternalServerError)
		return
	}

	// rip from-here
	res := &data.AuthResponse{
		Scope:     "openid profile email",
		ExpiresIn: 86400,
		TokenType: "Bearer",
	}
	privateKey, err := jwk.GetSigningKey()
	if err != nil {
		logger.ErrorContext(ctx, "Failed to encode JWKs", "error", err)
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
		return
	}

	err = authCode.ConsumeAuthCodeRequest(ctx, connection)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to consume auth code", "error", err)
		http.Error(w, "Failed to consume auth code", http.StatusInternalServerError)
		return
	}

	idToken := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"foo": "idToken",
	})
	if token, err := idToken.SignedString(privateKey); err != nil {
		logger.ErrorContext(ctx, "Failed to encode JWKs", "error", err)
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
		logger.ErrorContext(ctx, "Failed to encode JWKs", "error", err)
		http.Error(w, "Failed to encode JWKs", http.StatusInternalServerError)
		return
	} else {
		res.AccessToken = &token
		accessToken = nil
	}
	defer func() {
		res.AccessToken = nil
	}()
	// to here

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", "200 OK")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to encode json", "error", err)
		http.Error(w, "Failed to encode json", http.StatusInternalServerError)
		return
	}
}

func (c *Controller) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	connection, err := c.Db.Acquire(ctx)
	defer connection.Release()
	if err != nil {
		logger.Error("Failed to acquire connection", "error", err)
		http.Error(w, "Failed to acquire connection", http.StatusInternalServerError)
		return
	}

	valid := c.validateRedirect(ctx, w, connection, r.URL.Query().Get("redirect_uri"))
	if !valid {
		return
	}

	// TODO implement following:
	//valid = c.validateScope(ctx, w, connection, r.URL.Query().Get("scope"))
	//valid = c.validateResponseType(ctx, w, connection, r.URL.Query().Get("response_type"))
	//valid = c.validateClientId(ctx, w, connection, r.URL.Query().Get("client_id"))
	//valid = c.validateAudience(ctx, w, connection, r.URL.Query().Get("audience"))

	c.Renderer.Render(r.Context(), w, "authorize.html", "Login")
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

	realmId := ctx.Value(middleware.RealmIDContextKey).(int)
	redir := r.URL.Query().Get("redirect_uri")
	valid := c.validateRedirect(ctx, w, connection, redir)
	if !valid {
		return
	}

	user := new(data.User)
	err = user.GetUser(ctx, connection, r.FormValue("username"), realmId)
	if err != nil {
		logger.Warn("User not found", "username", r.FormValue("username"), "error", err)
		c.Renderer.RenderErr(ctx, w, "authorize.html", "Login", "Invalid Password or Username")
		return
	}
	if user.Password != r.FormValue("password") {
		logger.WarnContext(ctx, "User's credentials did not match!", "username", r.FormValue("username"))
		c.Renderer.RenderErr(ctx, w, "authorize.html", "Login", "Invalid Password or Username")
		return
	}
	authReq := new(data.AuthCodeRequest)
	err = authReq.CreateAuthCodeRequest(ctx, connection, user.ID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to create auth code", "error", err)
		c.Renderer.RenderErr(ctx, w, "authorize.html", "Login", "Auth code error, please try again later")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%s?code=%s", redir, authReq.AuthCode), http.StatusSeeOther) // hehe, stupid shit going down right here ;)
}

func (c *Controller) GetLandingPage(w http.ResponseWriter, r *http.Request) {
	if s := r.URL.Path; s != "/" { // make sure that the shit does not effect other pages.
		http.Redirect(w, r, "/error/404", http.StatusSeeOther)
		return
	}
	c.Renderer.Render(r.Context(), w, "authorize.html", "Login")
}
func (c *Controller) GetRegisterPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	connection, err := c.Db.Acquire(ctx)
	defer connection.Release()
	if err != nil {
		logger.Error("Failed to acquire connection", "error", err)
		http.Error(w, "Failed to acquire connection", http.StatusInternalServerError)
		return
	}
	countries, err := data.GetCountries(ctx, connection)
	if err != nil {
		logger.Error("Failed to acquire connection", "error", err)
		http.Error(w, "Failed to acquire connection", http.StatusInternalServerError)
		return
	}
	c.Renderer.RenderWithData(ctx, w, "register.html", utils.Page{
		RealmName:          r.Context().Value(middleware.RealmContextKey).(string),
		EnableRegistration: true,
		PageMeta:           utils.MetaData{PageTitle: "Login"},
		Countries:          countries,
	})
}

func (c *Controller) validateRedirect(ctx context.Context, w http.ResponseWriter, connection data.DbActions, redir string) bool {
	logger := ctx.Value(middleware.LoggerContextKey).(*slog.Logger)
	valid, err := data.CheckRedirectURI(ctx, connection, redir)
	if redir == "" || err != nil || !valid {
		logger.ErrorContext(ctx, "Invalid redirect_uri", "redirect_uri", redir)
		c.Renderer.RenderWithData(ctx, w, "error.html", utils.Page{
			PageMeta:   utils.MetaData{PageTitle: "Bad Request"},
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid redirect_uri",
		})
		return false
	}
	return true
}

func (c *Controller) GetAdminPage(w http.ResponseWriter, r *http.Request) {
	c.Renderer.RenderWithData(r.Context(), w, "admin_login.html", utils.Page{PageMeta: utils.MetaData{PageTitle: "Admin Login"}})
}
func (c *Controller) ErrorRedirect(w http.ResponseWriter, r *http.Request) {
	status := parseStatusCode(r.PathValue("status"))
	c.Renderer.RenderWithData(r.Context(), w, "error.html", utils.Page{
		PageMeta:   utils.MetaData{PageTitle: "Error"},
		StatusCode: status,
		Error:      http.StatusText(status),
	})
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
	c.Renderer.RenderWithData(r.Context(), w, "user_settings.html", utils.Page{PageMeta: utils.MetaData{PageTitle: "User Settings"}})
}
func (c *Controller) GetAdminDashboardPage(w http.ResponseWriter, r *http.Request) {
	c.Renderer.RenderWithData(r.Context(), w, "admin_panel.html", utils.Page{PageMeta: utils.MetaData{PageTitle: "Admin Dashboard"}})
}
