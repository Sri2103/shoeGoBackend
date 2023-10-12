package productService

import (
	productModel "github.com/sri2103/shoeMart/internal/app/product/models"
)

type ProductService interface {
	AddProduct(product *productModel.Product) (*productModel.Product, error)
	FindById(id string) (*productModel.Product, error)
	GetAll() ([]*productModel.Product, error)
}
