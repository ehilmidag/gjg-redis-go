package models

type UserCreateDTO struct {
	UserID      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	Points      float64 `json:"points"`
	Rank        int64   `json:"rank"`
}

type SignIn struct {
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	CountryCode string `json:"country_code"`
}

type UserCreateEntity struct {
	UserID         string  `json:"user_id"`
	DisplayName    string  `json:"display_name"`
	HashedPassword string  `json:"hashed_password"`
	Points         float64 `json:"points"`
	Rank           int64   `json:"rank"`
	CountryCode    string  `json:"country_code"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
}
