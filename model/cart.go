package model

type Cart struct {
	SkuId       string        `json:"sku_id"`
	Num         int           `json:"num"`
	Id          int64         `json:"id"`
	Specs       string        `json:"specs"`
	Type        int           `json:"type"`
	Name        string        `json:"name"`
	Point       int           `json:"point"`
	MainImage   string        `json:"mainImage"`
	Status      bool          `json:"status"`
	Stock       int           `json:"stock"`
	Description []interface{} `json:"description"`
	Sort        int           `json:"sort"`
}
