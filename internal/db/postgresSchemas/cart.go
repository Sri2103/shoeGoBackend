package postgresModels

import (
	"github.com/google/uuid"
	cartModel "github.com/sri2103/shoeMart/internal/app/cart/model"
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primaryKey"`
	CartID    uuid.UUID
	ProductID uuid.UUID 
	Quantity  int    
	Cart      Cart  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Product   Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Cart struct {
	gorm.Model
	ID     uuid.UUID `gorm:"primaryKey"`
	UserID uuid.UUID
	Status bool
	User   User  
}

type CartAggregate struct {
	Cart      *Cart
	CartItems []*CartItem
}

func (c *CartItem) ToEntity() *cartModel.CartItem {
	return &cartModel.CartItem{
		ID:        c.ID.String(),
		ProductID: c.ProductID.String(),
		Quantity:  c.Quantity,
		Product: *c.Product.ToEntity(),
	}
}

func (c *CartItem) FromEntity(item *cartModel.CartItem) error {
	if len(item.ID) == 0 {
		c.ID = uuid.New()
	} else {
		id, err := uuid.Parse(item.ID)
		c.ID = id
		if err != nil {
			return err
		}
	}
	productID, err := uuid.Parse(item.ProductID)

	if err != nil {
		return err
	}
	c.ProductID = productID
	c.Quantity = item.Quantity

	return nil

}

func (c *CartAggregate) ToEntity() (*cartModel.Cart, error) {

	cartItems := make([]cartModel.CartItem, len(c.CartItems))

	for i, p := range c.CartItems {
		mp := p.ToEntity()

		cartItems[i] = *mp
	}

	var cart = &cartModel.Cart{
		ID:       c.Cart.ID.String(),
		UserID:   c.Cart.UserID.String(),
		Status:   c.Cart.Status,
		CarItems: cartItems,
	}

	return cart, nil
}

func (c *CartAggregate) FromEntity(cart *cartModel.Cart) error {
	if len(cart.ID) == 0 {
		c.Cart.ID = uuid.New()
	} else {
		id, err := uuid.Parse(cart.ID)
		c.Cart.ID = id
		if err != nil {
			return err
		}
	}

	userId, err := uuid.Parse(cart.UserID)

	if err != nil {
		return err
	}

	c.Cart.UserID = userId

	c.Cart.Status = cart.Status

	items := make([]*CartItem, len(cart.CarItems))

	for i, p := range cart.CarItems {
		item := &p
		items[i] = &CartItem{}
		items[i].FromEntity(item)
		items[i].CartID = c.Cart.ID
	}
	c.CartItems = items
	return nil
}
