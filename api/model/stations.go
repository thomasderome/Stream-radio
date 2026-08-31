package model

type StationsData struct {
	Name  string `json:"name"`
	Url   string `json:"url_resolved"`
	Image string `json:"favicon"`
	Tags  string `json:"tags"`
}

type Stations struct {
	Name  string `db:"name"`
	Url   string `db:"url"`
	Image string `db:"image"`
}

type Tags struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}
