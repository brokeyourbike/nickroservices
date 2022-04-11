package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/brokeyourbike/nickroservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		rx := regexp.MustCompile(`/([0-9]+)`)
		g := rx.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(g[0][1])
		if err != nil {
			http.Error(w, "Invalid ID value", http.StatusBadRequest)
		}

		p.l.Println("got ID", id)

		p.updateProduct(id, w, r)

		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	productList := data.GetProducts()

	err := productList.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	product := &data.Product{}

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", product)

	data.AddProduct(product)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	product := &data.Product{}

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Unable to update product", http.StatusBadRequest)
		return
	}
}
