package productsHandler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	productRepo "github.com/sri2103/shoeMart/internal/app/product/repository"
	productService "github.com/sri2103/shoeMart/internal/app/product/service"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

func SetupProductRoutes(productRepo productRepo.ProductRepo, router *mux.Router, config *utils.Config, logger hclog.Logger) {
	productService := productService.NewProductService(productRepo, logger, config)
	productHandler := NewProduct(productService)
	
	productRouter := router.PathPrefix("/products").Subrouter()
	
	
	productRouter.HandleFunc("/testGetRoute", productHandler.TestRoute).Methods(http.MethodGet)
	productRouter.HandleFunc("/testPostRoute", productHandler.TestPostRoute).Methods(http.MethodPost)
	productRouter.HandleFunc("/upload", productHandler.AddProduct).Methods(http.MethodPost)
	productRouter.HandleFunc("/allProducts", productHandler.GetAll).Methods(http.MethodGet)
	productRouter.HandleFunc("/{id}", productHandler.FindProduct).Methods(http.MethodGet)

}
