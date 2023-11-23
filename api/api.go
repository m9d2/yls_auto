package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	token := c.Query("token")
	fmt.Println(token)
	targetTimeStr := c.Query("time")
	orderStr := c.Query("data")
	delayTime := c.Query("delayTime")

	num, err := strconv.Atoi(delayTime)
	if err != nil {
		JSON(c, err)
		return
	}
	if num < 200 {
		JSON(c, errors.New("time "))
		return
	}

	order := model.NewOrder()
	err = json.Unmarshal([]byte(orderStr), &order)
	if err != nil {
		JSON(c, err)
		return
	}
	if err != nil {
		JSON(c, err)
		return
	}

	now := time.Now()

	targetTime, err := time.Parse("15:04:05", targetTimeStr)
	if err != nil {
		JSON(c, err)
		return
	}

	nextTargetTime := time.Date(now.Year(), now.Month(), now.Day(), targetTime.Hour(), targetTime.Minute(), targetTime.Second(), 0, now.Location())
	if now.After(nextTargetTime) {
		nextTargetTime = nextTargetTime.Add(24 * time.Hour)
	}
	timeUntilNextTargetTime := nextTargetTime.Sub(now)

	timer := time.NewTimer(timeUntilNextTargetTime)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	roundTimer := time.NewTicker(time.Duration(num * 1000 * 1000))

	defer timer.Stop()
	defer roundTimer.Stop()
	for {
		select {
		case <-timer.C:
			for {
				select {
				case <-roundTimer.C:
					err := service.OrderSubmit(*order, token)
					if err != nil {
						err = sse.Encode(c.Writer, sse.Event{
							Event: "message",
							Data:  err.Error(),
						})
						if err != nil {
							return
						}
					} else {
						err = sse.Encode(c.Writer, sse.Event{
							Event: "message",
							Data:  "兑换成功！",
						})
						if err != nil {
							return
						}
					}
					flusher.Flush()
				case <-c.Writer.CloseNotify():
					return
				}
			}
		}
	}
}
