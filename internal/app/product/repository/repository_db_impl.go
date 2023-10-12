package productRepo

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	productModel "github.com/sri2103/shoeMart/internal/app/product/models"
	"github.com/sri2103/shoeMart/internal/app/utils"
	postgresModels "github.com/sri2103/shoeMart/internal/db/postgresSchemas"
	"gorm.io/gorm"
)

type PostgresDBimpl struct {
	db     *gorm.DB
	config *utils.Config
	logger hclog.Logger
}

func NewPostgresDBimpl(db *gorm.DB, config *utils.Config, logger hclog.Logger) *PostgresDBimpl {
	return &PostgresDBimpl{db: db, config: config, logger: logger}
}

func (r *PostgresDBimpl) AddProduct(product *productModel.Product) (*productModel.Product, error) {
	var item = new(postgresModels.Product)

	err := item.FromEntity(product)

	if err != nil {
		return nil, err
	}
	err = r.db.Create(item).Error
	if err != nil {
		return nil, err
	}
	product.ID = item.ID.String()
	return product, err
}

func (r *PostgresDBimpl) FindById(id string) (*productModel.Product, error) {
	Id, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var pT postgresModels.Product
	if err = r.db.Where("id=?", Id).First(&pT).Error; err != nil {
		return nil, err
	}

	return pT.ToEntity(), nil
}

func (r *PostgresDBimpl) GetAll() ([]*productModel.Product, error) {

	var postgresProducts []postgresModels.Product

	if err := r.db.Find(&postgresProducts).Error; err != nil {

		return nil, err
	}
	var modelProduct = make([]*productModel.Product, len(postgresProducts))
	for i, p := range postgresProducts {
		mp := p.ToEntity()

		modelProduct[i] = mp
	}
	return modelProduct, nil
}
