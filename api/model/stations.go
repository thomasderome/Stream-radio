package model

type StationsData struct {
	Name  string `json:"name"`
	Url   string `json:"url_resolved"`
	Image string `json:"favicon"`
	Tags  string `json:"tags"`
}

type PlayersData struct {
	IsPlaying bool    `db:"play" json:"is_playing"`
	Volume    float64 `db:"volume" json:"volume"`

	StationId   int    `db:"id" json:"station_id"`
	StationName string `db:"name" json:"station_name"`
	StationUrl  string `db:"url" json:"-"`
	StationImg  string `db:"img" json:"station_img"`
}
