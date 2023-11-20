package model

type Good struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	List []struct {
		Name             string        `json:"name"`
		Point            int           `json:"point"`
		SortPoint        int           `json:"sortPoint"`
		MainImage        string        `json:"mainImage"`
		SkuId            string        `json:"sku_id"`
		Stock            int           `json:"stock"`
		IsSubscribe      int           `json:"is_subscribe"`
		IsDiscount       bool          `json:"is_discount"`
		DiscountPoint    int           `json:"discountPoint"`
		IsSaleTime       bool          `json:"is_saleTime"`
		GoodsCategory    string        `json:"goodsCategory"`
		Description      []interface{} `json:"description"`
		ShopWindowImages []string      `json:"shopWindowImages"`
	} `json:"list"`
}
