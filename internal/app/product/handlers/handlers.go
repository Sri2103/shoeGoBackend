package productsHandler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	productModel "github.com/sri2103/shoeMart/internal/app/product/models"
	productService "github.com/sri2103/shoeMart/internal/app/product/service"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

type Product struct {
	productService productService.ProductService
}

func NewProduct(productService productService.ProductService) *Product {
	return &Product{
		productService: productService,
	}
}

func (p *Product) TestRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Test Response added"))
}
func (p *Product) TestPostRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Test Response Post response added"))
}

func (p *Product) AddProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Add product to DB")
	var item = new(productModel.Product)
	err := utils.FromJson(item, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = item.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := p.productService.AddProduct(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	err = utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "Add product",
		Data:    result,
	}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (p *Product) FindProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	item, err := p.productService.FindById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "Fetch product",
		Data:    item,
	}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (p *Product) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := p.productService.GetAll()

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "Fetch all products",
		Data:    products,
	}, w)

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
