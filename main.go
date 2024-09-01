package main

import (
	"authmantle-sso/handlers"
	"authmantle-sso/middleware"
	"authmantle-sso/oidc"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mainMiddleware := middleware.RegisterMiddlewares(
		middleware.RequestLogging,
		middleware.RenderTemplateContext,
	)
	cM := middleware.RegisterMiddlewares(middleware.EnsureSession)
	adnM := middleware.RegisterMiddlewares(
		middleware.AdminLock,
	)

	openRouter := http.NewServeMux()
	adminRouter := http.NewServeMux()
	adminOpenRouter := http.NewServeMux()
	closedRouter := http.NewServeMux()

	// OIDC registration
	openRouter.HandleFunc("GET /.well-known/openid-configuration", oidc.HandleWellKnown)
	for _, v := range oidc.ConfiguredRoutes {
		openRouter.HandleFunc(fmt.Sprintf("%s %s", v.Method, v.Endpoint), v.FunctionPTR)
	}
	//openRouter.HandleFunc("GET /.well-known/jwks.json", oidc.HandleJWKs)
	//openRouter.HandleFunc("POST /authorize", oidc.HandleAuth)
	//openRouter.HandleFunc("POST /oauth/token.json", oidc.HandleNewToken)
	// openRouter.HandleFunc("POST /oauth/refresh.json", oidc.HandleRefreshToken)
	// openRouter.HandleFunc("POST /oauth/revoke.json", oidc.HandleRevocation)
	// openRouter.HandleFunc("POST /oauth/logout.json", oidc.HandleLogout)

	// UI registration
	openRouter.HandleFunc("GET /register", handlers.GetRegisterPage)
	openRouter.HandleFunc("GET /authorize", handlers.GetLoginPage)

	openRouter.HandleFunc("GET /error/{status}", handlers.ErrorRedirect)
	openRouter.HandleFunc("GET /", handlers.GetLandingPage)

	adminRouter.HandleFunc("GET /", handlers.GetAdminDashboardPage)
	adminOpenRouter.HandleFunc("GET /", handlers.GetAdminPage)
	adminOpenRouter.HandleFunc("POST /", handlers.AdminLogin)

	closedRouter.HandleFunc("GET /", handlers.GetUserSettings)

	router := http.NewServeMux()
	router.Handle("/v1/", http.StripPrefix("/v1", openRouter))
	router.Handle("/protected/", http.StripPrefix("/protected", cM(closedRouter)))
	router.Handle("/adm_console/", http.StripPrefix("/adm_console", adnM(adminRouter)))
	router.Handle("/adm_login/", http.StripPrefix("/adm_login", adminOpenRouter))
	router.Handle("GET /admin/console", http.RedirectHandler("/adm_console/", http.StatusSeeOther))

	srv := &http.Server{
		Addr:    "localhost:8443",
		Handler: mainMiddleware(router),
	}
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)
	go func() {
		fmt.Println("Server started at localhost 8443")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Fatal error: %v\n", err)
			os.Exit(1)
		}
	}()
	<-sigint
	fmt.Println("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v", err)
		os.Exit(1)
	}
}