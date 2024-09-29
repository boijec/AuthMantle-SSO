package main

import (
	"authmantle-sso/controllers"
	"authmantle-sso/data"
	"authmantle-sso/middleware"
	"authmantle-sso/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO reformat to not have the entire universe in main - rule #1 of the weekend-warrior: don't be horny in main!
func main() {
	dbConnection, err := data.InitializePool()
	defer dbConnection.Close()
	if err != nil {
		slog.Error("Failed to initialize database connection", "error", err)
		os.Exit(1)
	}
	err = dbConnection.Ping()
	if err != nil {
		slog.Error("Failed to ping database post startup", "error", err)
		os.Exit(1)
	}
	renderer, err := utils.InitializeTemplates()
	if err != nil {
		slog.Error("Failed to initialize templates", "error", err)
		os.Exit(1)
	}
	uuid.EnableRandPool()
	slog.SetLogLoggerLevel(slog.LevelInfo)
	mainMiddleware := middleware.RegisterMiddlewares(
		middleware.RequestLogging,
	)
	rm := middleware.RealmMiddleware{
		Db: &dbConnection,
	}
	realmMiddleware := middleware.RegisterMiddlewares(
		rm.EnsureRealm,
	)

	controller := controllers.Controller{
		Db:       &dbConnection,
		Renderer: &renderer,
		BaseUrl:  "http://localhost:8443",
	}

	// ConfiguredRoutes global map of configured routes for OIDC discovery
	ConfiguredRoutes := map[string]*controllers.EndpointHelper{
		"JWKsUri":               {"GET", "/{realm}/.well-known/jwks.json", controller.HandleJWKs},
		"AuthorizationEndpoint": {"POST", "/{realm}/authorize", controller.HandleAuth},
		"TokenEndpoint":         {"POST", "/{realm}/oauth/token.json", controller.HandleNewToken},
	}

	controller.Discovery = &ConfiguredRoutes

	openRouter := http.NewServeMux()
	realmRouter := http.NewServeMux()
	//closedRouter := http.NewServeMux()
	router := http.NewServeMux()

	router.Handle("/", middleware.GetRoute(controller.GetLandingPage))
	router.Handle("/error/{status}", middleware.GetRoute(controller.ErrorRedirect))

	// OIDC registration
	realmRouter.HandleFunc("GET /{realm}/.well-known/openid-configuration", controller.HandleWellKnown)
	for _, v := range ConfiguredRoutes {
		realmRouter.HandleFunc(fmt.Sprintf("%s %s", v.Method, v.Endpoint), v.FunctionPTR)
	}
	//openRouter.HandleFunc("GET /.well-known/jwks.json", oidc.HandleJWKs)
	//openRouter.HandleFunc("POST /authorize", oidc.HandleAuth)
	//openRouter.HandleFunc("POST /oauth/token.json", oidc.HandleNewToken)
	// openRouter.HandleFunc("POST /oauth/refresh.json", oidc.HandleRefreshToken)
	// openRouter.HandleFunc("POST /oauth/revoke.json", oidc.HandleRevocation)
	// openRouter.HandleFunc("POST /oauth/logout.json", oidc.HandleLogout)

	// UI registration
	realmRouter.HandleFunc("GET /{realm}/register", controller.GetRegisterPage)
	realmRouter.HandleFunc("GET /{realm}/authorize", controller.GetLoginPage)

	//closedRouter.HandleFunc("GET /user-info", controller.GetUserSettings)

	router.Handle("/v1/", http.StripPrefix("/v1", openRouter))
	router.Handle("/v1/oidc/", http.StripPrefix("/v1/oidc", realmMiddleware(realmRouter)))

	srv := &http.Server{
		Addr:    "localhost:8443", // TODO export to env
		Handler: mainMiddleware(router),
	}
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)
	go func() {
		slog.Info("Server started at localhost 8443")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Fatal error", "error", err)
			os.Exit(1)
		}
	}()
	<-sigint
	slog.Info("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v", err)
		os.Exit(1)
	}
}
