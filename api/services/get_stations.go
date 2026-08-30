package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http"
	"radio_stream/model"
	"radio_stream/utils"
	"strings"
)

func Init_get_stations() {
	run()
}

func run() {
	stations := get_list_stations()
	process_stations(stations)
}

func get_list_stations() *[]model.StationsData {
	resp, err := http.Get("https://de1.api.radio-browser.info/json/stations/bycountrycodeexact/fr")
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

func process_stations(data *[]model.StationsData) {
	tags_set := make(utils.Set)
	langs_set := make(utils.Set)

	for _, station := range *data {
		split_and_add(station.Tags, tags_set)
		split_and_add(station.Langs, langs_set)
	}

	query, args := query_maker(tags_set, "tags")
	_, err := DB.Exec(query, args...)
	if err != nil {
		log.Fatal("Error inserting tags: ", err)
	}
}

func query_maker(data utils.Set, table string) (string, []any) {
	args := make([]any, 0, len(data))
	placeholders := make([]string, 0, len(data))

	for key, _ := range data {
		placeholders = append(placeholders, "(?)")
		args = append(args, key)
	}

	return fmt.Sprintf("INSERT OR IGNORE INTO %s (name) VALUES %s", table, strings.Join(placeholders, ",")), args
}

func split_and_add(s string, set utils.Set) {
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
