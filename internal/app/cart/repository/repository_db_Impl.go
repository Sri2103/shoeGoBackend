package cartRepo

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	cartModel "github.com/sri2103/shoeMart/internal/app/cart/model"
	"github.com/sri2103/shoeMart/internal/app/utils"
	postgresModels "github.com/sri2103/shoeMart/internal/db/postgresSchemas"
	"gorm.io/gorm"
)

type CartRepositoryImpl struct {
	db     *gorm.DB
	config *utils.Config
	logger hclog.Logger
}

func NewCartRepo(db *gorm.DB, config *utils.Config, logger hclog.Logger) *CartRepositoryImpl {
	return &CartRepositoryImpl{
		db,
		config,
		logger,
	}
}

func (r *CartRepositoryImpl) AddToCart(cart *cartModel.Cart) (*cartModel.Cart, error) {
	var agg = &postgresModels.CartAggregate{
		Cart: &postgresModels.Cart{},
	}

	err := agg.FromEntity(cart)

	if err != nil {
		r.logger.Error("Unable convert from entity", err)
		return nil, err
	}

	err = r.db.Create(agg.Cart).Error

	if err != nil {
		r.logger.Error("Unable convert from entity", err)
		return nil, err
	}

	err = r.db.Create(agg.CartItems).Error

	if err != nil {
		r.logger.Error("Unable convert from entity", err)
		return nil, err
	}

	return agg.ToEntity()
	// return nil, nil
}
func (r *CartRepositoryImpl) FetchCart(userId string) (*cartModel.Cart, error) {

	cart := new(postgresModels.Cart)
	uid, err := uuid.Parse(userId)
	if err != nil {
		r.logger.Error("Parsing to uuid", "error", err)
		return nil, err
	}
	err = r.db.First(cart, "user_id=?", uid).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart.ID = uuid.New()
			cart.UserID = uid

			err = r.db.Create(cart).Error
			if err != nil {
				r.logger.Error("could not create new cart", "error", err)
				return nil, err
			}
		} else {
			r.logger.Error("Fetch cart failed", "error", err)
			return nil, err
		}

	}
	var cartItems []*postgresModels.CartItem

	err = r.db.Preload("Product").Where("cart_id=?", cart.ID).Find(&cartItems).Error

	if err != nil {
		r.logger.Error("Fetch cartItems failed", "error", err)
		return nil, err
	}

	cartTotal := &postgresModels.CartAggregate{
		Cart:      cart,
		CartItems: cartItems,
	}

	c, err := cartTotal.ToEntity()

	return c, err
}
func (r *CartRepositoryImpl) UpdateCart(cartId string, cartItem *cartModel.CartItem) error {
	var Item = new(postgresModels.CartItem)

	err := r.db.First(Item, "cart_id = ? AND product_id = ?", cartId, cartItem.ProductID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Item.ID = uuid.New()
			Item.CartID, _ = uuid.Parse(cartId)
			Item.ProductID, _ = uuid.Parse(cartItem.ProductID)
			Item.Quantity = cartItem.Quantity

			err = r.db.Create(&Item).Error

			if err != nil {
				r.logger.Error("Adding new cartItem failed", "error", err)
				return err
			}
		} else {
			r.logger.Error("Fetch cartItem failed", "error", err)
			return err
		}
	}

	r.logger.Debug("%+v", Item)

	Item.Quantity = cartItem.Quantity

	err = r.db.Save(&Item).Error

	if err != nil {
		r.logger.Error("Update cartItem failed", "error", err)
		return err
	}

	return nil
}
