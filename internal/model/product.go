package model

type ProductRespone struct {
	ProductId     string   `json:"id"`
	Name          string   `json:"name"`
	Price         int64    `json:"price"`
	ImageUrl      string   `json:"imageUrl"`
	Stock         int16    `json:"stock"`
	Condition     string   `json:"condition"`
	Tags          []string `json:"tags"`
	IsPurchasable bool     `json:"isPurchasable"`
	PurchaseCount int8     `json:"purchaseCount"`
}

type ProductRequest struct {
	Name          string   `json:"name" validate:"required,min=5,max=60"`
	Price         int64    `json:"price" validate:"required,min=0"`
	ImageUrl      string   `json:"imageUrl" validate:"required,url"`
	Stock         int16    `json:"stock"`
	Condition     string   `json:"condition" validate:"required,oneof=new second"`
	Tags          []string `json:"tags" validate:"required"`
	IsPurchasable bool     `json:"isPurchasable" validate:"required"`
}

type ProductFilter struct {
	Condition      *string   `json:"condition"`
	Keyword        *string   `json:"keyword"`
	SortBy         *string   `json:"sortBy"`
	OrderBy        *string   `json:"orderBy"`
	MaxPrice       *int      `json:"maxPrice"`
	MinPrice       *int      `json:"minPrice"`
	Tags           *[]string `json:"tags"`
	UserOnly       *bool     `json:"userOnly"`
	ShowEmptyStock *bool     `json:"showEmptyStock"`
	Limit          *int      `json:"limit"`
	Offset         *int      `json:"offset"`
}

type StockRequest struct {
	Stock int16 `json:"stock" validate:"min=0"`
}

type BuyRequest struct {
	BankAccountId        string `json:"bankAccountId" validate:"required"`
	ProductId            string `json:"productId" validate:"required"`
	PaymentProofImageUrl string `json:"paymentProofImageUrl" validate:"required,url"`
	Quantity             int16  `json:"quantity" validate:"required,min=1"`
}
