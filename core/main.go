package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mqnoy/logistics-app/core/config"
	"github.com/mqnoy/logistics-app/core/handler"
	"github.com/mqnoy/logistics-app/core/middleware"
	"github.com/mqnoy/logistics-app/core/model"
	"gorm.io/gorm"

	_goodHttpDelivery "github.com/mqnoy/logistics-app/core/good/delivery"
	_goodRepoMySQL "github.com/mqnoy/logistics-app/core/good/repository/mysql"
	_godUseCase "github.com/mqnoy/logistics-app/core/good/usecase"
	transaction "github.com/mqnoy/logistics-app/core/transaction_manager/repository"
)

var (
	appCfg config.Configuration
)

type AppCtx struct {
	mysqlDB *gorm.DB
}

func init() {
	appCfg = config.AppConfig
}

func main() {
	appCtx := AppCtx{
		mysqlDB: config.InitMySQLDatabase(appCfg),
	}

	// Auto migration
	if appCfg.MigrateConfig.AutoMigrate {
		err := appCtx.mysqlDB.AutoMigrate(
			&model.Good{},
			&model.GoodStock{},
		)

		if err != nil {
			log.Println(err.Error())
		}
	}

	// The HTTP Server
	addr := appCfg.Server.Address()
	server := &http.Server{
		Addr:    addr,
		Handler: AppHandler(appCtx),
	}

	log.Printf("server running on %s\n", addr)

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

	// Initialize trx manager
	txManager := transaction.New(appctx.mysqlDB)

	// Initialize Repository
	goodRepoMySQL := _goodRepoMySQL.New(appctx.mysqlDB)

	// Initialize UseCase
	goodUseCase := _godUseCase.New(goodRepoMySQL)

	// Fallback
	mux.NotFound(handler.FallbackHandler)

	// Initialize handler
	_goodHttpDelivery.New(mux, goodUseCase)

	// Print all routes
	chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return mux
}
