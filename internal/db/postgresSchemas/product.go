package postgresModels

import (
	"github.com/google/uuid"
	productModel "github.com/sri2103/shoeMart/internal/app/product/models"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string
	Description string
	Image       string
	Price       float32
}

func (p *Product) ToEntity() *productModel.Product {
	return &productModel.Product{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Image:       p.Image,
		Price:       p.Price,
	}
}
func (p *Product) FromEntity(item *productModel.Product) error {

	if len(item.ID) == 0 {
		p.ID = uuid.New()
	} else {
		id, err := uuid.Parse(item.ID)
		p.ID = id
		if err != nil {
			return err
		}
	}

	p.Name = item.Name
	p.Description = item.Description
	p.Image = item.Image
	p.Price = item.Price

	return nil
}

