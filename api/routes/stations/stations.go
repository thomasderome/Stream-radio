package stations

import (
	"fmt"
	"net/http"
	model "radio_stream/model/routes_model"
	s "radio_stream/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	group := engine.Group("/stations")

	group.GET("/favorite", getFavoritesStation)
	group.GET("/station", getStationWithParams)
}

func getFavoritesStation(c *gin.Context) {
	favorite := make([]model.Station, 0)
	err := s.DB.Select(&favorite, "SELECT stations.id, stations.name, stations.img, GROUP_CONCAT(DISTINCT tags.name) AS tags FROM stations JOIN favorite ON stations.id = favorite.station_id JOIN station_tags ON station_tags.station_id = stations.id JOIN tags ON tags.id = station_tags.tag_id ORDER BY favorite.listened_at")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range favorite {
		favorite[i].Tags = strings.Split(favorite[i].TagsOutDB, ",")
	}

	c.JSON(http.StatusOK, favorite)
}

func getStationWithParams(c *gin.Context) {
	var query model.StationGetQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var args []any
	var request strings.Builder
	request.WriteString("SELECT stations.id, stations.name, stations.img, COALESCE(GROUP_CONCAT(DISTINCT tags.name), '') AS tags FROM stations LEFT JOIN station_tags ON station_tags.station_id = stations.id LEFT JOIN tags ON tags.id = station_tags.tag_id")

	if query.Search != "" {
		request.WriteString(" WHERE stations.name LIKE ?")
		args = append(args, query.Search+"%")
	}
	if query.Tags != "" {
		if query.Search == "" {
			request.WriteString(" WHERE ")
		} else {
			request.WriteString(" AND ")
		}

		parts := strings.Split(query.Tags, ",")

		for i, split := range parts {
			request.WriteString("? = station_tags.tag_id")
			if i != len(parts)-1 {
				request.WriteString(" AND ")
			}

			args = append(args, split)
		}
	}

	request.WriteString(" GROUP BY stations.name LIMIT 50 OFFSET ?")
	args = append(args, query.Offset)
	stations := make([]model.Station, 0)

	err := s.DB.Select(&stations, request.String(), args...)
	fmt.Print(request.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range stations {
		stations[i].Tags = strings.Split(stations[i].TagsOutDB, ",")
	}

	c.JSON(http.StatusOK, stations)
}
