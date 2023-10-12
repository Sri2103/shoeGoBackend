package productRepo

type ProductRepoImpl struct {
}

func NewProductRepoImpl() *ProductRepoImpl {
	return &ProductRepoImpl{}
}
func (p *ProductRepoImpl) AddProduct() error {
	return nil
}
