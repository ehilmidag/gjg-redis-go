package models

type LeaderBoard struct {
	DisplayName string  `json:"display_name"`
	Point       float64 `json:"point"`
	Rank        int64   `json:"rank"`
	Country     string  `json:"country"`
}
