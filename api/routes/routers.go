package routes

import "github.com/gin-gonic/gin"
import (
	"radio_stream/routes/players"
)

func Register_all_route(router *gin.Engine) {
	players.Register(router)
}
