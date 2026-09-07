package model

type SetPlayStateRequest struct {
	PlayState bool `json:"play_state"`
}

type SetVolumeRequest struct {
	Volume float64 `json:"volume"`
}
