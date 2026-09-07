package players

import (
	"net/http"
	"radio_stream/model"
	"radio_stream/utils"

	"github.com/gin-gonic/gin"

	s "radio_stream/services"
)

func Register(engine *gin.Engine) {
	r := engine.Group("/players")

	r.PUT("/set_play_state", setPlayState)
	r.PUT("/set_volume", setVolume)

	r.GET("/get_play_data", getPlayState)
}

func setPlayState(c *gin.Context) {
	var body model.SetPlayStateRequest

	if err := utils.BodyBinder(c, &body); err != nil {
		return
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

func setVolume(c *gin.Context) {
	var body model.SetVolumeRequest
	if err := utils.BodyBinder(c, &body); err != nil {
		return
	}

	volume, err := s.PlayersService.SetVolume(body.Volume)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	if body.Volume > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Volume must be >= 1"})
	}

	c.JSON(http.StatusOK, gin.H{"volume": volume})
}

func getPlayState(c *gin.Context) {
	data, err := s.GetStationData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, data)
}
