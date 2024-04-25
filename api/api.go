package api

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"

	"github.com/enzujp/konoha-mhi/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type API struct {
	Server *http.Server
	Config	*config.Configuration
}


func (a *API) Serve() error {
	a.Server = &http.Server{
		Addr: fmt.Sprintf(":%d", a.Config.Port),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1024 * 1024,
		Handler: a.CustomHandler(),
	}

	fmt.Printf("Server running on Port :%v\n", a.Config.Port)
	return a.Server.ListenAndServe()
}

func (a *API) CustomHandler() http.Handler{
	router := chi.NewRouter()
	router.Use(
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders: []string{"Link"},
			AllowCredentials: false,
			MaxAge: 300,
		}),	
	)
	router.Get("/", Home)
	router.Mount("/users", a.UserRoutes(router))
	return router
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Annyeonghaseyo, welcome to the homepage!"))
}


func decodeRequestBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}