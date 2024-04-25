package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)


func (a *API) UserRoutes(router chi.Router) http.Handler {
	userRouter := chi.NewRouter()
	userRouter.Post("/{userID}/disburse", a.DisburseFunds)
	userRouter.Get("/details/{userID}", a.GetUserDetails)
	userRouter.Get("/{transactionID}", a.GetTransactionDetails)
	return userRouter
}