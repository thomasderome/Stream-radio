package main

import (
	"radio_stream/middleware"
	"radio_stream/routes"
	"radio_stream/services"

	"github.com/gin-gonic/gin"
)

func main() {
	services.InitServices()

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	routes.Register_all_route(router)

	err := router.Run(":3000")
	if err != nil {
		panic(err)
	}
}
