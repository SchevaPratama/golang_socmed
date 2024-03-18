package converter

import (
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
)

func ProductConverter(product *entity.Product) *model.ProductRespone {
	return &model.ProductRespone{
		ProductId:     product.ID,
		Name:          product.Name,
		Price:         product.Price,
		ImageUrl:      product.ImageUrl,
		Stock:         product.Stock,
		Condition:     product.Condition,
		Tags:          product.Tags,
		IsPurchasable: product.IsPurchasable,
		PurchaseCount: product.PurchaseCount,
	}
}
