package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"ysl_auto/model"
)

const (
	baseUrl     = "https://wemember.platform-loreal.cn"
	contentType = "application/json"
	appHash     = "ysl"
	originType  = "miniApp"
	userAgent   = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.42(0x18002a31) NetType/WIFI Language/zh_CN"
)

func Goods(token string) ([]model.Good, error) {
	url := "/api/getGoodsList"
	params := `{"type": 1,"apiListId": 3,"goods_num": 200}`
	var goods []model.Good
	response, err := request(http.MethodPost, url, params, token)
	if err != nil {
		return goods, err
	}
	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return goods, err
	}

	err = json.Unmarshal([]byte(jsonData), &goods)
	if err != nil {
		return goods, err
	}
	return goods, nil
}

func AddCart(skuId string, num int, token string) error {
	url := "/api/exchange-gift/cart-add"
	params := fmt.Sprintf(`{"num": %d,"source": 2,"sku_id": %s}`, num, skuId)
	_, err := request(http.MethodPost, url, params, token)
	if err != nil {
		return err
	}
	return nil
}

//	{
//		"product": [{
//			"sku_id": "UQhf3VTJisRbHjejiB9FFc",
//			"num": 1
//		}]
//	}
func CheckCart(skuId string, num int, token string) error {
	url := "/api/exchange-gift/cart-check"
	params := fmt.Sprintf(`{"num": %d,"source": 2,"sku_id": %s}`, num, skuId)
	_, err := request(http.MethodPost, url, params, token)
	if err != nil {
		return err
	}
	return nil
}

//	{
//		"address_id": 68463539002624,
//		"is_cart": 0,
//		"is_subscribe": 1,
//		"product": [{
//			"sku_id": "UQhf3VTJisRbHjejiB9FFc",
//			"num": 1
//		}]
//	}
func OrderSubmit(entity model.Order, token string) error {
	url := "/api/exchange-gift/ordinary/order-submit"
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	_, err = request(http.MethodPost, url, string(jsonData), token)
	if err != nil {
		return err
	}
	return nil
}

func CartList(token string) (carts []model.Cart, err error) {
	url := "/api/exchange-gift/cart-list"
	_, err = request(http.MethodGet, url, "", token)
	if err != nil {
		return carts, err
	}
	return carts, nil
}

func AddressList(token string) (model.Address, error) {
	url := "/api/exchange-gift/address-list?page_size=10&page=1"
	var address model.Address
	response, err := request(http.MethodGet, url, "", token)
	if err != nil {
		return address, err
	}
	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return address, err
	}

	err = json.Unmarshal([]byte(jsonData), &address)
	return address, nil
}

func GetUser(token string) (model.User, error) {
	url := "/api/user"
	var user model.User
	response, err := request(http.MethodGet, url, "", token)
	if err != nil {
		return user, err
	}
	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal([]byte(jsonData), &user)
	return user, nil
}

func request(method string, url string, params string, token string) (model.Response, error) {
	var response model.Response
	url = baseUrl + url
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(params)))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("App-Hash", appHash)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin-Type", originType)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	if response.HttpCode != 200 {
		slog.Error("request fail", "url", url, "code", response.HttpCode, "message", response.Message)
		return response, errors.New(response.Message)
	}
	if response.StatusCode != 0 {
		slog.Error("operate fail", "url", url, "code", response.StatusCode, "message", response.Message)
		return response, errors.New(response.Message)
	}
	return response, nil
}
