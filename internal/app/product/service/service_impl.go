package productService

import (
	"github.com/hashicorp/go-hclog"
	productModel "github.com/sri2103/shoeMart/internal/app/product/models"
	productRepo "github.com/sri2103/shoeMart/internal/app/product/repository"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

type ProductServiceImpl struct {
	productRepo productRepo.ProductRepo
	logger      hclog.Logger
	config      *utils.Config
}

func NewProductService(repo productRepo.ProductRepo, logger hclog.Logger, config *utils.Config) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo: repo,
		logger:      logger,
		config:      config,
	}
}

func (service *ProductServiceImpl) AddProduct(product *productModel.Product) (*productModel.Product, error) {
	item, err := service.productRepo.AddProduct(product)
	return item, err
}

func (service *ProductServiceImpl) FindById(id string) (*productModel.Product, error) {
	product, err := service.productRepo.FindById(id)
	return product, err
}

func (service *ProductServiceImpl) GetAll() ([]*productModel.Product, error) {
	products, err := service.productRepo.GetAll()

	return products, err
}
