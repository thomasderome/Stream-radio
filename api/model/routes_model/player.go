package model_routes

type SetPlayStateRequest struct {
	PlayState bool `json:"play_state"`
}

type SetVolumeRequest struct {
	Volume float64 `json:"volume"`
}

type SetStationsRequest struct {
	StationId string `json:"station_id"`
}
