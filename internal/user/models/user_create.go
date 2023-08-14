package models

type UserCreateDTO struct {
	UserID      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	Points      float64 `json:"points"`
}

type UserResponseModel struct {
	UserID      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	Point       float64 `json:"point"`
	Rank        int64   `json:"rank"`
}

type SignIn struct {
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	Country     string `json:"country"`
}

type UserCreateEntity struct {
	UserID         string  `json:"user_id"`
	DisplayName    string  `json:"display_name"`
	HashedPassword string  `json:"hashed_password"`
	Country        string  `json:"country"`
	Points         float64 `json:"points"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
}
