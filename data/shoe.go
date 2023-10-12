package data
// todo: reformat

import (
	"encoding/json"
	"fmt"
	"io"
)

type Shoe struct {
	ID    int     `json:"id"`
	Brand string  `json:"brand,omitempty"`
	Size  float64 `json:"size" binding:"required"` // size in meters
	Model string  `json:"model"`
}

var ErrorShoeNotFound = fmt.Errorf("Shoe not found")

type Shoes []*Shoe

var shoes = []*Shoe{
	{
		ID:    1,
		Brand: "Adidas",
		Size:  1.2,
		Model: "Running",
	},
	{
		ID:    3,
		Brand: "Puma",
		Size:  1.63,
		Model: "Soccer",
	},
}

func GetShoes() Shoes {
	return shoes
}
func (p *Shoe) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Shoe) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func FindShoe(id int) (*Shoe, int, error) {
	for i, p := range shoes {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrorShoeNotFound

}
