package models

type SendScore struct {
	UserID      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	Country     string  `json:"country"`
	ScoreWorth  float64 `json:"score"`
}

type SendScoreDto struct {
	UserID     string  `json:"user_id"`
	ScoreWorth float64 `json:"score_worth"`
	TimeStamp  int64   `json:"time_stamp"`
}

type SendScoreEntity struct {
	UserID     string  `json:"user_id"`
	TotalScore float64 `json:"total_score"`
	TimeStamp  int64   `json:"time_stamp"`
}
