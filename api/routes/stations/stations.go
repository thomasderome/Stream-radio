package stations

import (
	"net/http"
	model "radio_stream/model/routes_model"
	s "radio_stream/services"
	"radio_stream/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	group := engine.Group("/stations")

	group.GET("/favorite", getFavoritesStation)
	group.POST("/favorite", addFavoriteStation)
	group.DELETE("/favorite/:id", removeFavoriteStation)
	group.GET("/station", getStationWithParams)
}

func getFavoritesStation(c *gin.Context) {
	favorite := make([]model.Station, 0)
	err := s.DB.Select(&favorite, "SELECT stations.id, stations.name, stations.img, COALESCE(GROUP_CONCAT(DISTINCT tags.name), '') AS tags, 1 AS is_favorite, favorite.volume FROM favorite JOIN stations ON stations.id = favorite.station_id LEFT JOIN station_tags ON station_tags.station_id = stations.id LEFT JOIN tags ON tags.id = station_tags.tag_id GROUP BY stations.id ORDER BY favorite.added_at DESC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range favorite {
		favorite[i].Tags = strings.Split(favorite[i].TagsOutDB, ",")
	}

	c.JSON(http.StatusOK, favorite)
}

func addFavoriteStation(c *gin.Context) {
	var body model.FavoriteRequest
	if err := utils.BodyBinder(c, &body); err != nil {
		return
	}

	var count int
	err := s.DB.Get(&count, "SELECT COUNT(*) FROM stations WHERE id = ?", body.StationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	_, err = s.DB.Exec("INSERT OR IGNORE INTO favorite (station_id) VALUES (?)", body.StationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"station_id": body.StationId, "is_favorite": true})
}

func removeFavoriteStation(c *gin.Context) {
	stationId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid station id"})
		return
	}

	_, err = s.DB.Exec("DELETE FROM favorite WHERE station_id = ?", stationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"station_id": stationId, "is_favorite": false})
}

func getStationWithParams(c *gin.Context) {
	var query model.StationGetQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var args []any
	var request strings.Builder
	request.WriteString("SELECT stations.id, stations.name, stations.img, COALESCE(GROUP_CONCAT(DISTINCT tags.name), '') AS tags, favorite.id IS NOT NULL AS is_favorite, favorite.volume FROM stations LEFT JOIN favorite ON favorite.station_id = stations.id LEFT JOIN station_tags ON station_tags.station_id = stations.id LEFT JOIN tags ON tags.id = station_tags.tag_id")

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

		request.WriteString(" tags.name IN (")
		for i := range parts {
			if i != 0 {
				request.WriteString(",")
			}
			request.WriteString("?")
		}
		request.WriteString(")")
		for _, split := range parts {
			args = append(args, split)
		}
	}

	request.WriteString(" GROUP BY stations.name LIMIT 50 OFFSET ?")
	args = append(args, query.Offset)
	stations := make([]model.Station, 0)

	err := s.DB.Select(&stations, request.String(), args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range stations {
		stations[i].Tags = strings.Split(stations[i].TagsOutDB, ",")
	}

	c.JSON(http.StatusOK, stations)
}
