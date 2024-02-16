package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func renderError(c echo.Context, err string) error {
	return renderErrorWithCode(c, http.StatusInternalServerError, err)
}

func renderErrorWithCode(c echo.Context, code int, err string) error {
	return c.JSON(code, map[string]string{"status": "error", "error": err})
}

func renderOK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "ok", "data": data})
}

type AuthHandler struct {
	RedirectURL string
}

func (h *AuthHandler) API(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Request().Header.Get("Remote-User")
		if user == "" {
			return c.JSON(
				http.StatusUnauthorized,
				map[string]string{
					"status":       "error",
					"error":        "unknown_user",
					"redirect_url": h.RedirectURL,
				},
			)
		}
		return next(c)
	}
}
