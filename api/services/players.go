package services

import "log"

type PlayersServiceStruct struct {
	is_playing bool
	url_steam  string
}

var PlayersService *PlayersServiceStruct

func InitPlayersService() {
	PlayersService = new(PlayersServiceStruct)
}

func (PlayersService *PlayersServiceStruct) SetUrlSteam(url_steam string) {
	PlayersService.url_steam = url_steam
	_, err := DB.Exec("UPDATE players SET url_steam = ?", url_steam)
	if err != nil {
		log.Fatal(err)
	}
}
