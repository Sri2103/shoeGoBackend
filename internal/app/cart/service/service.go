package cartService

import cartModel "github.com/sri2103/shoeMart/internal/app/cart/model"

type CartService interface {
	AddToCart(cart *cartModel.Cart) (*cartModel.Cart, error)
	FetchCart(userId string) (*cartModel.Cart, error)
	UpdateCart(cartId string,cartItem *cartModel.CartItem) (error)
}
