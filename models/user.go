package models

type UserGameMessage struct {
	GameMode         string `json:"game_mode"`
	AreaCode         int64  `json:"area_code"`
	CurrentLoginTime int64  `json:"current_login_time"`
}

type User struct {
	UserGame UserGameMessage
	UserId   int64 `json:"user_id"`
}

type CurrentPlayerSocketMessage struct {
	AreaCode         int64 `json:"area_code"`
	CurrentLoginTime int64 `json:"current_login_time"`
}

type CurrentPlayer struct {
	CurrentPlayerSocketMessage CurrentPlayerSocketMessage
	UserId                     int64 `json:"user_id"`
}
