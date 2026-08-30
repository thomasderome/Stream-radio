package main

import (
	"radio_stream/services"

	"github.com/gin-gonic/gin"
)

func main() {
	services.InitDB()

	router := gin.Default()
	err := router.Run(":3000")
	if err != nil {
		panic(err)
	}
}
