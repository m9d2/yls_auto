package main

import (
	"ysl_auto/router"
)

func main() {
	r := router.InitRouter()
	err := r.Run(":8081")
	if err != nil {
		return
	}
}
