package data
// todo: reformat

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// product struct defines the data of the product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func GetProducts() Products {
	return ProductList
}

var ProductList = []*Product{
	{
		ID:          1,
		Name:        "IPhone",
		Description: "14promax",
		Price:       1705.12,
		SKU:         "ghj456",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	{
		ID:          2,
		Name:        "IPhone",
		Description: "14promax",
		Price:       1705.12,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1
}

var ErrorProductNotFound = fmt.Errorf("Product not found")

func (p *Products) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func getNextIds() int {
	lp := ProductList[len(ProductList)-1]
	return lp.ID + 1
}

func AddProduct(p *Product) {
	p.ID = getNextIds()

	ProductList = append(ProductList, p)
}

func FindProduct(id int) (*Product, int, error) {
	for i, p := range ProductList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrorProductNotFound

}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := FindProduct(p.ID)
	if err != nil {
		return err
	}
	p.ID = id

	ProductList[pos] = p
	return nil
}
