package main

import (
	"main/routes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	routes.Use(router)
	http.ListenAndServe(":5000", router)
}
