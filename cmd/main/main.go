package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	cartHandlers "github.com/sri2103/shoeMart/internal/app/cart/handlers"
	cartRepo "github.com/sri2103/shoeMart/internal/app/cart/repository"
	productsHandler "github.com/sri2103/shoeMart/internal/app/product/handlers"
	productRepo "github.com/sri2103/shoeMart/internal/app/product/repository"
	userHandler "github.com/sri2103/shoeMart/internal/app/user/handlers"
	userRepository "github.com/sri2103/shoeMart/internal/app/user/repository"
	"github.com/sri2103/shoeMart/internal/app/utils"
	"github.com/sri2103/shoeMart/internal/db"
)

func main() {

	hcLogger := utils.NewLogger()

	config := utils.NewConfig(hcLogger)

	validator := utils.NewValidation()

	// logging := log.New(os.Stdout, "", log.LstdFlags)
	r := SetUp(config, hcLogger, validator)
	headersOk := gohandlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := gohandlers.AllowedOrigins([]string{"*"})
	methodsOk := gohandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	m := gohandlers.LoggingHandler(os.Stdout, r)

	ch := gohandlers.CORS(originsOk, headersOk, methodsOk)
	srv := &http.Server{
		Handler:      ch(m),
		Addr:        fmt.Sprintf(":%s", config.ServerPort),
		ErrorLog:     hcLogger.StandardLogger(&hclog.StandardLoggerOptions{}),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		hcLogger.Info("Server starting on :", config.ServerPort)
		if err := srv.ListenAndServe(); err != nil {
			hcLogger.Error(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c

	log.Printf("Got signal %v\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

}

func SetUp(config *utils.Config, logger hclog.Logger, validator *utils.Validation) *mux.Router {

	database, err := db.ConnectToDB(config)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.MigrateModels(database)
	if err != nil {
		log.Fatal("Failed to perform auto-migration:", err)
	}

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	productRepoImpl := productRepo.NewPostgresDBimpl(database, config, logger)
	productsHandler.SetupProductRoutes(productRepoImpl, r, config, logger)

	userRepoImpl := userRepository.NewUserRepo(database, logger)
	userHandler.SetupUserRoutes(userRepoImpl,r,config,logger,validator)

	cartRepoImpl := cartRepo.NewCartRepo(database, config, logger)
	cartHandlers.SetupCartRoutes(cartRepoImpl,r,config,logger,validator)


	return r
}
