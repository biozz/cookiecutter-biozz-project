package server

import (
	"github.com/labstack/echo/v4"
)

func (h *Web) ListItems(c echo.Context) error {
	items, err := h.Storage.GetItems(
		c.Request().Context(),
	)
	if err != nil {
		return renderError(c, err.Error())
	}
	return renderOK(c, items)
}
