package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"radio_stream/model"
	"radio_stream/utils"
	"strings"
)

var COUNTRIES = []string{"FR", "MC"}

func InitGetStations() {
	stations := make([]model.StationsData, 0)
	for _, country := range COUNTRIES {
		stations = append(stations, getListStations(country)...)
	}

	log.Printf("%d stations récupérées", len(stations))
	processStations(stations)
}

func getListStations(countryCode string) []model.StationsData {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://de1.api.radio-browser.info/json/stations/search?countrycode="+countryCode+"&order=clickcount&reverse=true&hidebroken=true&limit=100000", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:152.0) Gecko/20100101 Firefox/152.0")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error get stations: ", err)
	}
	defer resp.Body.Close()

	databyte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error decoding stations to byte: ", err)
	}

	var stations []model.StationsData
	err = json.Unmarshal(databyte, &stations)
	if err != nil {
		log.Fatal("Error decoding stations data: ", err)
	}

	return stations
}

func processStations(data []model.StationsData) {
	tx := DB.MustBegin()
	defer tx.Rollback()

	tx.MustExec("DELETE FROM station_tags")
	tx.MustExec("DELETE FROM tags")

	tagsSet := make(utils.Set)

	for _, station := range data {
		splitAndAdd(station.Tags, tagsSet)
	}

	if len(tagsSet) > 0 {
		query, args := QueryMakerInsertOneColumn("tags", "name", tagsSet)
		_, err := tx.Exec(query, args...)
		if err != nil {
			log.Fatal("Error inserting tags: ", err)
		}
	}

	idTags := make(map[string]int64, len(tagsSet))
	rows, err := tx.Query("SELECT id, name FROM tags")

	if err != nil {
		log.Fatal("Error getting tags: ", err)
	}

	for rows.Next() {
		var id int64
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error scanning rows: ", err)
		}

		idTags[name] = id
	}

	stationsSet := make(utils.Set, len(data))
	urlsSet := make(utils.Set, len(data))
	for _, station := range data {
		nameTrimmed := strings.TrimSpace(station.Name)
		if nameTrimmed == "" || station.Url == "" || stationsSet.Contains(nameTrimmed) || urlsSet.Contains(station.Url) {
			continue
		}

		stationsSet.Add(nameTrimmed)
		urlsSet.Add(station.Url)

		result := tx.QueryRow("INSERT INTO stations(name, url, img) VALUES(?,?,?) ON CONFLICT(name) DO UPDATE SET img = excluded.img RETURNING id", nameTrimmed, station.Url, station.Image)
		var id int64
		err = result.Scan(&id)
		if err != nil {
			log.Printf("Station ignorée (%s): %s", nameTrimmed, err)
			continue
		}

		for tag := range strings.SplitSeq(station.Tags, ",") {
			tagId, ok := idTags[cleanTag(tag)]
			if !ok {
				continue
			}

			_, err = tx.Exec("INSERT OR IGNORE INTO station_tags (tag_id, station_id) VALUES(?, ?)", tagId, id)
			if err != nil {
				log.Fatal("Error link tag on stations: ", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Error commit stations: ", err)
	}
}

func cleanTag(tag string) string {
	return strings.ToLower(strings.TrimSpace(tag))
}

func splitAndAdd(s string, set utils.Set) {
	if s == "" {
		return
	}

	for value := range strings.SplitSeq(s, ",") {
		tag := cleanTag(value)
		if tag == "" {
			continue
		}
		set.Add(tag)
	}
}
