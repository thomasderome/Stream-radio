package model_routes

type Station struct {
	Id        int      `db:"id" json:"station_id"`
	Name      string   `db:"name" json:"name"`
	Image     string   `db:"img" json:"favicon"`
	TagsOutDB string   `db:"tags" json:"-"`
	Tags      []string `json:"tags"`

	IsFavorite bool     `db:"is_favorite" json:"is_favorite"`
	Volume     *float64 `db:"volume" json:"volume,omitempty"`
}

type FavoriteRequest struct {
	StationId int `json:"station_id"`
}

type StationGetQuery struct {
	Offset int    `form:"offset"`
	Search string `form:"search"`
	Tags   string `form:"tags"`
}
