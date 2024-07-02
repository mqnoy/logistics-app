package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mqnoy/logistics-app/core/handler"
	"github.com/mqnoy/logistics-app/core/middleware"
	"gorm.io/gorm"
)

var (
	Address = "0.0.0.0:8080"
)

type AppCtx struct {
	mysqlDB *gorm.DB
}

func main() {
	appCtx := AppCtx{}

	// The HTTP Server
	server := &http.Server{
		Addr:    Address,
		Handler: AppHandler(appCtx),
	}

	log.Println("server started..")

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func AppHandler(appctx AppCtx) http.Handler {
	mux := chi.NewRouter()

	// Setup middleware
	mux.Use(chiMiddleware.RealIP)
	mux.Use(middleware.PanicRecoverer)

	// Fallback
	mux.NotFound(handler.FallbackHandler)

	// TODO: Initialize handler
	mux.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.ParseResponse(w, r, "hc", map[string]interface{}{
				"result": time.Now().Unix(),
			}, nil)
		})
	})

	// Print all routes
	chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return mux
}
