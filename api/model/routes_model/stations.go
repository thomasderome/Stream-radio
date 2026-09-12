package model_routes

type Station struct {
	Id        int      `db:"id" json:"station_id"`
	Name      string   `db:"name" json:"name"`
	Image     string   `db:"img" json:"favicon"`
	TagsOutDB string   `db:"tags" json:"-"`
	Tags      []string `json:"tags"`
}

type StationGetQuery struct {
	Offset int    `form:"offset"`
	Search string `form:"search"`
	Tags   string `form:"tags"`
}
