package model

type StationsData struct {
	Name  string `json:"name"`
	Url   string `json:"url_resolved"`
	Image string `json:"favicon"`
	Tags  string `json:"tags"`
}

type PlayersData struct {
	IsPlaying bool    `db:"play"`
	Volume    float64 `db:"volume"`

	StationName string `db:"name"`
	StationUrl  string `db:"url"`
	StationImg  string `db:"img"`
}
