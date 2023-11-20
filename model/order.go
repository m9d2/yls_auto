package model

type Order struct {
	AddressId   int64 `json:"address_id"`
	IsCart      int   `json:"is_cart"`
	IsSubscribe int   `json:"is_subscribe"`
	Product     []struct {
		SkuId string `json:"sku_id"`
		Num   int    `json:"num"`
	} `json:"product"`
}

func NewOrder() *Order {
	return &Order{
		AddressId:   0,
		IsCart:      0,
		IsSubscribe: 1,
		Product:     nil,
	}
}
