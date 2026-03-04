package types

import "time"

type Product struct {
	ID			int			`json:"product_id"`
	Name 		string		`json:"product_name"`
	URL 		string		`json:"url"`
	CreatedAt 	time.Time	`json:"created_at"`
	Prices		[]PriceData `json:"prices"`
}

type PriceData struct {
	Source		string		`json:"source"`
	Price		float64		`json:"price"`
	InStock		bool		`json:"in_stock"`
	Timestamp	time.Time	`json:"timestamp"`
}

type UserProduct struct {
	ProductID 		int 		`json:"product_id"`
	ProductName		string		`json:"product_name"`
	ImageUrl		string		`json:"image_url"`
	AddedAt			time.Time	`json:"added_at"`
	LastCheckedAt 	time.Time	`json:"last_checked_at"`
	LowestPrice		float64		`json:"lowest_price"`
	LowestSource	string 		`json:"lowest_source"`
	InStock			bool		`json:"in_stock"`
}