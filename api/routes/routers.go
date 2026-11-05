package routes

import "github.com/gin-gonic/gin"
import (
	"radio_stream/routes/player"
	"radio_stream/routes/stations"
)

func Register_all_route(router *gin.Engine) {
	player.Register(router)
	stations.Register(router)
}
