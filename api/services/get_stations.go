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

func Init_get_stations() {
	run()
}

func run() {
	stations := getListStations()
	processStations(stations)
}

func getListStations() *[]model.StationsData {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://de1.api.radio-browser.info/json/stations/search?countrycode=FR&order=clickcount&reverse=true&hidebroken=true", nil)
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

	var stations *[]model.StationsData
	err = json.Unmarshal(databyte, &stations)
	if err != nil {
		log.Fatal("Error decoding stations data: ", err)
	}

	return stations
}

func processStations(data *[]model.StationsData) {
	tagsSet := make(utils.Set)

	for _, station := range *data {
		splitAndAdd(station.Tags, tagsSet)
	}

	query, args := QueryMakerInsertOneColumn("tags", "name", tagsSet)
	_, err := DB.Exec(query, args...)
	if err != nil {
		log.Fatal("Error inserting tags: ", err)
	}

	idTags := make(map[string]int64, len(tagsSet))
	rows, err := DB.Query("SELECT id, name FROM tags")

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

	stationsSet := make(utils.Set, len(*data))
	for _, station := range *data {
		if stationsSet.Contains(station.Name) {
			continue
		}
		stationsSet.Add(strings.TrimSpace(station.Name))

		result := DB.QueryRow("INSERT INTO stations(name, url, img) VALUES(?,?,?) ON CONFLICT(name) DO UPDATE SET id = id RETURNING id", station.Name, station.Url, station.Image)
		var id int64
		err = result.Scan(&id)
		if err != nil {
			log.Fatal("Impossible get id of stations insert: ", err)
		}

		if station.Tags == "" {
			continue
		}
		for tag := range strings.SplitSeq(station.Tags, ",") {
			_, err = DB.Exec("INSERT INTO station_tags (tag_id, station_id) VALUES(?, ?)", idTags[tag], id)
			if err != nil {
				log.Fatal("Error link tag on stations: ", err)
			}
		}
	}
}

func splitAndAdd(s string, set utils.Set) {
	if s == "" {
		return
	}

	for value := range strings.SplitSeq(s, ",") {
		trim := strings.TrimSpace(value)
		if set.Contains(strings.TrimSpace(trim)) {
			continue
		}
		set.Add(trim)
	}
}
