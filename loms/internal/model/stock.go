package model

type Stock struct {
	SKU           SKUType `json:"sku"`
	TotalCount    uint32  `json:"total_count"`
	ReservedCount uint32  `json:"reserved"`
}
