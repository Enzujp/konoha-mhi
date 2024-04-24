package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)


func (a *API) UserRoutes(router chi.Router) http.Handler {
	router.Get("/{id}", a.GetUserDetails)
	router.Post("/{id}/disburse", a.DisburseFunds)
	return router
}