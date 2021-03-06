package pivnet

import (
	"fmt"
	"net/http"

	"github.com/pivotal-cf/go-pivnet/logger"
)

type ProductsService struct {
	client Client
	l      logger.Logger
}

type Product struct {
	ID   int    `json:"id,omitempty" yaml:"id,omitempty"`
	Slug string `json:"slug,omitempty" yaml:"slug,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ProductsResponse struct {
	Products []Product `json:"products,omitempty"`
}

func (p ProductsService) List() ([]Product, error) {
	url := "/products"

	var response ProductsResponse
	_, _, err := p.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return []Product{}, err
	}

	return response.Products, nil
}

func (p ProductsService) Get(slug string) (Product, error) {
	url := fmt.Sprintf("/products/%s", slug)

	var response Product
	_, _, err := p.client.MakeRequest(
		"GET",
		url,
		http.StatusOK,
		nil,
		&response,
	)
	if err != nil {
		return Product{}, err
	}

	return response, nil
}
