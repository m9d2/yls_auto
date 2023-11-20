package model

type Address struct {
	List []struct {
		Id             int64  `json:"id"`
		DeliveryName   string `json:"delivery_name"`
		DeliveryMobile string `json:"delivery_mobile"`
		Province       string `json:"province"`
		City           string `json:"city"`
		District       string `json:"district"`
		Address        string `json:"address"`
		IsDefault      int    `json:"is_default"`
	} `json:"list"`
}
