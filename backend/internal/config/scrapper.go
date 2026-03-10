package config

import "time"


type SiteConfig struct {
	Name 			string
	Domain 			string
	SearchURLFormat string
	RateLimit 		time.Duration
	Selectors 		ScraperSelectors
	IDExtractor 	func(string) string
}

// contains css selectors for extracting product data
type ScraperSelectors struct {
	ProductContainer 	string
	ProductName 		string
	PriceSelectors		[]string // multiple price selectors for fallback
	ProductImage		string
	URLAttribute		string
	URLSelector			string
}

var AmazonConfig = SiteConfig{
	Name: 				"amazon",
	Domain: 			"www.amazon.com",
	SearchURLFormat:	"https://www.amazon.com/s?k=%s",
	RateLimit: 			5 * time.Second,
	Selectors: 	ScraperSelectors{
		ProductContainer: 	"div[data-component-type's-search-result']",
		ProductName: 		"h2 a span",
		PriceSelectors: []string{
			"span.a-price > span.a-offscreen", // Clean price for screen readers
			"span.a-price-whole", 			   // fallback: whole number part
			"span.a-price", 				   // fallback: any price element
		},
		ProductImage: 		"img.s-image",
		URLSelector: 		"h2 a",
		URLAttribute: 		"href",
	},
	IDExtractor: ext,
}

var WalmartConfig = SiteConfig{
	Name: 				"walmart",
	Domain: 			"www.walmart.com",
	SearchURLFormat:	"https://www.walmart.com/search?q=%s",
	RateLimit: 			5 * time.Second,
	Selectors: 	ScraperSelectors{
		ProductContainer: 	"div[data-item-id']",
		ProductName: 		"h2 a span",
		PriceSelectors: []string{
			"span.a-price > span.a-offscreen", // Clean price for screen readers
			"span.a-price-whole", 			   // fallback: whole number part
			"span.a-price", 				   // fallback: any price element
		},
		ProductImage: 		"img.s-image",
		URLSelector: 		"h2 a",
		URLAttribute: 		"href",
	},
	IDExtractor: ext,
}