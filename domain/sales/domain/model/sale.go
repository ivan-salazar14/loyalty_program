package model

type Sale struct {
	Id         string  `json:"sale_id,omitempty"`
	ProductId  string  `json:"product_id"`
	CustomerId string  `json:"user_id"`
	Price      float64 `json:"price"`
}
