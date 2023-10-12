package productModel

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float32 `json:"price" validate:"required"`
}


func (p *Product) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	err := validate.Struct(*p)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("field:", err.Field(), err.Tag())

		}
	}
	return err
}
