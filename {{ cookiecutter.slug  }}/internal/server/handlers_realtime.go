package server

import (
	"fmt"
	"net/http"
)

func (h *Web) broadcastItems(r *http.Request, data map[string]interface{}) error {
	clientID := r.Header.Get("X-Client-ID")
	err := h.Clients.Centrifugo.BroadcastAllExcept(fmt.Sprintf("items#%s", clientID), map[string]interface{}{"hello": "world"})
	return err
}
