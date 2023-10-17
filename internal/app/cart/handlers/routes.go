package cartHandlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"

	cartRepo "github.com/sri2103/shoeMart/internal/app/cart/repository"
	cartService "github.com/sri2103/shoeMart/internal/app/cart/service"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

func SetupCartRoutes(cartRepo cartRepo.CartRepository, router *mux.Router, config *utils.Config, logger hclog.Logger, validation *utils.Validation) {
	cartService := cartService.NewCartService(cartRepo, logger, config)
	cartHandlers := NewCart(cartService, logger, config, validation)

	cartRouter := router.PathPrefix("/cart").Subrouter()
	cartRouter.HandleFunc("/temp", cartHandlers.TemporaryFunc).Methods(http.MethodGet)
	cartRouter.HandleFunc("/add", cartHandlers.AddToCart).Methods(http.MethodPost)
	cartRouter.HandleFunc("/get/{userId}", cartHandlers.FetchCart).Methods(http.MethodGet)
	cartRouter.HandleFunc("/update", cartHandlers.UpdateCart).Methods(http.MethodPost)
	cartRouter.HandleFunc("/delete", cartHandlers.DeleteCartItem).Methods(http.MethodDelete)

}
