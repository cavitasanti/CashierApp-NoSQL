package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
)

type ProductRepository struct {
	db db.DB
}

func NewProductRepository(db db.DB) ProductRepository {
	return ProductRepository{db}
}

func (u *ProductRepository) ReadProducts() ([]model.Product, error) {
	// return []model.Product{}, nil // TODO: replace this
	load, _ := u.db.Load("products")
	var products []model.Product
	err := json.Unmarshal([]byte(load), &products)
	if err != nil {
		return nil, err
	}
	return products, nil

}

func (u *ProductRepository) ResetProducts() error {
	err := u.db.Reset("products", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}
