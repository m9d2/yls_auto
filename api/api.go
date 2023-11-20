package api

import (
	"github.com/gin-gonic/gin"
	"ysl_auto/model"
	"ysl_auto/service"
)

type WebHandler struct {
}

func (h WebHandler) InitRouter(g *gin.RouterGroup) {
	g.GET("goods", h.list)
	g.GET("address", h.address)
	g.GET("user", h.user)
	g.POST("exchange", h.exchange)
}

func (h WebHandler) list(c *gin.Context) {
	token := c.Query("token")
	goods, err := service.Goods(token)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, goods)
	}
}

func (h WebHandler) address(c *gin.Context) {
	token := c.Query("token")
	address, err := service.AddressList(token)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, address)
	}
}

func (h WebHandler) exchange(c *gin.Context) {
	token := c.Query("token")
	order := model.NewOrder()
	err := c.ShouldBind(&order)
	if err != nil {
		JSON(c, err)
		return
	}
	err = service.OrderSubmit(*order, token)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, nil)
	}
}

func (h WebHandler) user(c *gin.Context) {
	token := c.Query("token")
	user, err := service.GetUser(token)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, user)
	}
}
