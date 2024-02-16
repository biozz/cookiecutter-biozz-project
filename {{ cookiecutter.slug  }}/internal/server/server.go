package server

import (
	"context"
	"fmt"
	"net/http"

	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/clients"
	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/repository"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/oauth2"
)

type Web struct {
	Dev         bool
	Storage     repository.Storage
	Echo        *echo.Echo
	BaseURL     string
	Secret      []byte
	UserInfoURL string
	AuthConfig  *oauth2.Config
	Clients     clients.Clients
}

type NewServerParams struct {
	OIDCProviderURL  string
	OIDCClientID     string
	OIDCClientSecret string
	JWTSecret        string
	BaseURL          string
	Dev              bool
	Storage          repository.Storage
	Clients          clients.Clients
}

func New(params NewServerParams) Web {
	e := echo.New()
	provider, err := oidc.NewProvider(context.Background(), params.OIDCProviderURL)
	if err != nil {
		panic(err)
	}
	redirectURL := fmt.Sprintf("%s/auth/callback", params.BaseURL)
	return Web{
		Dev:     params.Dev,
		Storage: params.Storage,
		Echo:    e,
		BaseURL: params.BaseURL,
		Secret:  []byte(params.JWTSecret),
		AuthConfig: &oauth2.Config{
			ClientID:     params.OIDCClientID,
			ClientSecret: params.OIDCClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			Endpoint:     provider.Endpoint(),
		},
		Clients: params.Clients,
	}
}

func (h *Web) Start(bind string) {
	h.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Remote-User", "X-Client-ID"},
		MaxAge:       360,
	}))
	h.Echo.Use(middleware.Logger())

	// prom := prometheus.NewPrometheus("wow", nil)
	// prom.Use(w.Echo)

	auth := h.Echo.Group("/auth")
	auth.GET("/login", h.Login)
	auth.GET("/logout", h.Logout)
	auth.GET("/callback", h.Callback)
	auth.GET("/realtime/token", h.RealtimeToken)

	api := h.Echo.Group("/api", h.AuthMiddlewareAPI)

	api.GET("/items", h.ListItems)

	h.Echo.GET("/*", echo.WrapHandler(http.StripPrefix("/", h.Static())), h.AuthMiddleware)

	h.Echo.Logger.Fatal(h.Echo.Start(bind))
}
