package cartService

import (
	"github.com/hashicorp/go-hclog"
	cartModel "github.com/sri2103/shoeMart/internal/app/cart/model"
	cartRepo "github.com/sri2103/shoeMart/internal/app/cart/repository"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

type CartServiceImpl struct {
	cartRepository cartRepo.CartRepository
	logger         hclog.Logger
	config         *utils.Config
}

func NewCartService(repo cartRepo.CartRepository, logger hclog.Logger, config *utils.Config) *CartServiceImpl {
	return &CartServiceImpl{
		repo,
		logger,
		config,
	}
}

func (service *CartServiceImpl) AddToCart(cart *cartModel.Cart) (*cartModel.Cart, error) {
	c, err := service.cartRepository.AddToCart(cart)
	if err != nil {
		service.logger.Error("error adding cart from repo")
		return nil, err
	}
	return c, nil
}

func (service *CartServiceImpl) FetchCart(userId string) (*cartModel.Cart, error) {
	c, err := service.cartRepository.FetchCart(userId)
	if err != nil {
		service.logger.Error("error  fetching cart from repo")
		return nil, err
	}
	return c, nil
}

func (service *CartServiceImpl) UpdateCart(cartId string, cartItem *cartModel.CartItem) error {
	err := service.cartRepository.UpdateCart(cartId, cartItem)
	if err != nil {
		service.logger.Error("error updating cart in repo", err)
		return err
	}
	return nil
}

func (service *CartServiceImpl) DeleteCartItem(cartId string, cartItemId string) error {
	err := service.cartRepository.DeleteCartItem(cartId, cartItemId)
	if err != nil {
		service.logger.Error("error deleting item from cart", err)
	}
	return err
}
