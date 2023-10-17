package cartRepo

import cartModel "github.com/sri2103/shoeMart/internal/app/cart/model"

type CartRepository interface {
	AddToCart(cart *cartModel.Cart) (*cartModel.Cart, error)
	FetchCart(userId string) (*cartModel.Cart, error)
	UpdateCart(cartId string, cartItem *cartModel.CartItem) ( error)
	DeleteCartItem(cartId string, cartItemId string) error
}
