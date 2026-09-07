package players

import (
	"net/http"
	"radio_stream/model"

	"github.com/gin-gonic/gin"

	s "radio_stream/services"
)

func Register(engine *gin.Engine) {
	r := engine.Group("/players")

	r.PUT("/set_play_state", setPlayState)
}

func setPlayState(c *gin.Context) {
	var body model.SetPlayStateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if body.PlayState {
		err := s.PlayersService.Resume()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"play_state": body.PlayState})
		return
	}

	err := s.PlayersService.Pause()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"play_state": body.PlayState})
}
