package service

import (
	"github.com/gohouse/gorose/v2"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-admin/model"
	"log"
)

type ProductService struct {
}

func NewProductServer() *ProductService {
	return &ProductService{}
}

func (p *ProductService) CreateProduct(product *model.Product) error {
	productEntity := model.NewProductModel()
	err := productEntity.CreateProduct(product)
	if err != nil {
		log.Printf("ProductEntity.CreateProduct, err : %v", err)
		return err
	}
	return nil
}

func (p *ProductService) GetProductList() ([]gorose.Data, error) {
	productEntity := model.NewProductModel()
	productList, err := productEntity.GetProductList()
	if err != nil {
		log.Printf("ProductEntity.CreateProduct, err : %v", err)
		return nil, err
	}
	return productList, nil
}
