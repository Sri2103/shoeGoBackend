package cartModel

import productModel "github.com/sri2103/shoeMart/internal/app/product/models"

type CartItem struct {
	ID        string               `json:"id,omitempty"`
	ProductID string               `json:"productId"`
	Quantity  int                  `json:"quantity"`
	Product   productModel.Product `json:"product"`
}

type Cart struct {
	ID       string     `json:"id"`
	UserID   string     `json:"userId"`
	Status   bool       `json:"status"`
	CarItems []CartItem `json:"cartItems,omitempty"`
}
