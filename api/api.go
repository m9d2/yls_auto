package api

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"ysl_auto/model"
	"ysl_auto/service"
)

var (
	ch = make(chan string)
)

type WebHandler struct {
}

func (h WebHandler) InitRouter(g *gin.RouterGroup) {
	g.GET("goods", h.list)
	g.GET("address", h.address)
	g.GET("user", h.user)
	g.POST("exchange", h.exchange)
	g.GET("exchange-auto", h.exchangeAuto)
	g.GET("sse", h.sse)
}

func (h WebHandler) list(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("logs") == nil {
		session.Set("logs", "queue")
		session.Save()
	}

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

func (h WebHandler) exchangeAuto(c *gin.Context) {
	sid, err := c.Cookie("sid")
	if err != nil {
		return
	}
	go func() {
		for range time.Tick(1 * time.Second) {
			ch <- sid + "@Hello"
		}
	}()

	JSON(c, "ok")
}

func (h WebHandler) sse(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	sid, err := c.Cookie("sid")
	if err != nil {
		return
	}

	go func() {
		for {
			msg := <-ch
			str := strings.Split(msg, "@")

			fmt.Println("Message received in N:", msg)

			if sid == str[0] {
				err := sse.Encode(c.Writer, sse.Event{
					Event: "log",
					Data:  str[1],
				})
				if err != nil {
					return
				}
				c.Writer.WriteString(str[1])
				c.Writer.Flush()
			}

		}
	}()
}
