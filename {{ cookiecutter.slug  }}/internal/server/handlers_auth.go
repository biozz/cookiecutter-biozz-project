package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (h *Web) Login(c echo.Context) error {
	// Authelia has strict requirement for the state, it has to be
	// at least 8 characters long
	url := h.AuthConfig.AuthCodeURL("statestate")
	return c.Redirect(http.StatusFound, url)
}

func (h *Web) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	token, err := h.AuthConfig.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.String(http.StatusBadRequest, "Error exchanging code for token")
	}

	userInfo := getUserFromToken(c.Request().Context(), token.AccessToken, h.AuthConfig, h.UserInfoURL)
	if userInfo == nil {
		c.Logger().Errorf("unable to get user from token %v", err)
		return c.Redirect(http.StatusFound, h.BaseURL)
	}

	expires := time.Now().Add(3 * 7 * 24 * time.Hour)
	claims := &Claims{
		Email:    userInfo.Email,
		Username: userInfo.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokString, err := tok.SignedString(h.Secret)
	if err != nil {
		return renderError(c, "unable to generate auth token")
	}

	cookie := &http.Cookie{
		Name:    "{{ cookiecutter.slug }}_token",
		Value:   tokString,
		Path:    "/",
		Expires: expires,
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, h.BaseURL)
}

func (h *Web) GetUser(c echo.Context) error {
	user := getAuthorizedUser(c, h.Secret)
	return renderOK(c, map[string]string{"email": user.Email})
}

func (h *Web) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:    "{{ cookiecutter.slug }}_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, h.BaseURL)
}

func (h *Web) AuthMiddlewareAPI(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// userInfo := getAuthorizedUser(c, h.Secret)
		// if userInfo == nil {
		// 	// return renderErrorWithCode(c, http.StatusUnauthorized, "unauthorized")
		// 	return c.JSON(
		// 		http.StatusUnauthorized,
		// 		map[string]string{
		// 			"status":       "error",
		// 			"error":        "unknown_user",
		// 			"redirect_url": h.BaseURL + "/auth/login",
		// 		},
		// 	)
		//
		// }
		userInfo := UserInfo{Email: "test@example.com"}
		c.Set("email", userInfo.Email)
		return next(c)
	}
}

func (h *Web) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// userInfo := getAuthorizedUser(c, h.Secret)
		// if userInfo == nil {
		// 	return c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		// }
		userInfo := UserInfo{Email: "test@example.com"}
		c.Set("email", userInfo.Email)
		return next(c)
	}
}

type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
}

type Claims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func getAuthorizedUser(c echo.Context, secret []byte) *UserInfo {
	cookie, err := c.Cookie("{{ cookiecutter.slug }}_token")
	if err != nil {
		return nil
	}
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		c.Logger().Warnf("unable to parse jwt: %s", err.Error())
		return nil
	}
	if !tok.Valid {
		return nil
	}
	return &UserInfo{Email: claims.Email, Username: claims.Username}
}

func getUserFromToken(ctx context.Context, token string, config *oauth2.Config, userInfoURL string) *UserInfo {
	t := &oauth2.Token{AccessToken: token}
	client := config.Client(ctx, t)
	resp, err := client.Get(userInfoURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err == nil {
		return &userInfo
	}
	return nil
}

func (h *Web) RealtimeToken(c echo.Context) error {
	// expires := time.Now().Add(3 * 7 * 24 * time.Hour)
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{})
	tokString, err := tok.SignedString(h.Secret)
	if err != nil {
		return renderError(c, "unable to generate auth token")
	}
	return renderOK(c, map[string]string{"token": tokString})
}
