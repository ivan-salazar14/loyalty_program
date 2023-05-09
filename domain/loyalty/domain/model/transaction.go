package model

import "time"

type Transaction struct {
	UserID        string    `json:"userId"`
	TransactionId string    `json:"TransactionId"`
	Type          string    `json:"type"`
	Date          time.Time `json:"date"`
	Product       string    `json:"product,omitempty"`
	Value         float64   `json:"value,omitempty"`
	Percentage    float64   `json:"percentage,omitempty"`
	Points        string    `json:"points"`
}
