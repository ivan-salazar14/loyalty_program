package model

type Product struct {
	ProductId   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
}
