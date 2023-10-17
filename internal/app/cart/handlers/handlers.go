package cartHandlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	cartModel "github.com/sri2103/shoeMart/internal/app/cart/model"
	cartService "github.com/sri2103/shoeMart/internal/app/cart/service"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

type Cart struct {
	cartService cartService.CartService
	logger      hclog.Logger
	config      *utils.Config
	validator   *utils.Validation
}

func NewCart(service cartService.CartService, logger hclog.Logger, config *utils.Config, validator *utils.Validation) *Cart {
	return &Cart{
		service,
		logger,
		config,
		validator,
	}
}
func (cart *Cart) TemporaryFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test cart routes"))
	w.WriteHeader(http.StatusOK)
}

func (cart *Cart) AddToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cartData = new(cartModel.Cart)
	err := utils.FromJson(&cartData, r.Body)

	if err != nil {
		cart.logger.Error("[AddToCart] Error while parsing request body: ", err)
		w.WriteHeader(400)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: err.Error(),
		}, w)

		return
	}

	storedData, err := cart.cartService.AddToCart(cartData)

	if err != nil {
		cart.logger.Error("[AddToCart] Error in adding to cart: ", err)
		w.WriteHeader(500)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: err.Error(),
		}, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "CartItem created",
		Data:    storedData,
	}, w)

}

func (cart *Cart) FetchCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := mux.Vars(r)["userId"]
	if len(userId) == 0 {
		cart.logger.Error("FetchCart - User Id is empty")
		w.WriteHeader(400)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: "User id cannot be null or blank.",
		}, w)
	}
	fullCart, err := cart.cartService.FetchCart(userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: "Find the cart",
		}, w)
		return
	}
	w.WriteHeader(200)
	utils.ToJson(&utils.GenericResponse{
		Success: false,
		Message: "Find the cart",
		Data:    fullCart,
	}, w)

}

func (cart *Cart) UpdateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type RequestBody struct {
		CartId   string             `json:"cartId"`
		CartItem cartModel.CartItem `json:"cartItem"`
	}
	var body RequestBody
	err := utils.FromJson(&body, r.Body)

	fmt.Println(body, "Input data")

	if err != nil {
		cart.logger.Error("[UpdateCart] Error while parsing request body: ", err)
		w.WriteHeader(400)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: err.Error(),
		}, w)
		return
	}

	cartId := body.CartId
	cartItem := body.CartItem
	fmt.Println(cartId, cartItem, "Request Params in Req.Body")
	err = cart.cartService.UpdateCart(cartId, &cartItem)

	if err != nil {
		cart.logger.Error("[AddToCart] Error while parsing request body: ", err)
		w.WriteHeader(500)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: err.Error(),
		}, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "CartItem updated",
		// Data:    bod,
	}, w)
}

func (cart *Cart) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	values := r.URL.Query()
	cartId := values["cartId"][0]
	cartItemId := values["cartItemId"][0]
	err := cart.cartService.DeleteCartItem(cartId, cartItemId)

	if err != nil {
		cart.logger.Error("[DeleteCartItem] Error while deleting item from cart: ", err)
		w.WriteHeader(500)

		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: "Could not delete the item from cart.",
		}, w)
		return
	}

	w.WriteHeader(200)
	utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "Item deleted successfully from cart.",
	}, w)

}
