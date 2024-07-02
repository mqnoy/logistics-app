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
	_orderHttpDelivery "github.com/mqnoy/logistics-app/core/order/delivery"
	_orderRepoMySQL "github.com/mqnoy/logistics-app/core/order/repository/mysql"
	_orderUseCase "github.com/mqnoy/logistics-app/core/order/usecase"

	transaction "github.com/mqnoy/logistics-app/core/transaction_manager/repository"

	_userHttpDelivery "github.com/mqnoy/logistics-app/core/user/delivery/http"
	_userRepoMySQL "github.com/mqnoy/logistics-app/core/user/repository/mysql"
	_userUsecase "github.com/mqnoy/logistics-app/core/user/usecase"
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
			&model.Order{},
			&model.User{},
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
	orderRepoMySQL := _orderRepoMySQL.New(appctx.mysqlDB)
	userRepoMySQL := _userRepoMySQL.New(appctx.mysqlDB)

	// Initialize UseCase
	goodUseCase := _godUseCase.New(txManager, goodRepoMySQL)
	orderUseCase := _orderUseCase.New(txManager, orderRepoMySQL, goodUseCase)
	userUseCase := _userUsecase.New(userRepoMySQL)

	// Fallback
	mux.NotFound(handler.FallbackHandler)

	// Initialize handler
	_goodHttpDelivery.New(mux, goodUseCase)
	_orderHttpDelivery.New(mux, orderUseCase)
	_userHttpDelivery.New(mux, userUseCase)

	// Print all routes
	chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return mux
}
