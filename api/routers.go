package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)


func (a *API) UserRoutes(router chi.Router) http.Handler {
	router.Get("/{userID}", a.GetUserDetails)
	router.Post("/{userID}/disburse", a.DisburseFunds)
	router.Get("/{transactionID}", a.GetTransactionDetails)
	return router
}