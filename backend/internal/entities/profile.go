package entities

type Profile struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  string  `json:"patronymic"`
	Description string  `json:"description"`
	UserLogin   string  `json:"user_login"`
	TG          string  `json:"tg"`
	PrevPoints  float64 `json:"prev_points"`
	CurrPoints  float64 `json:"curr_points"`
	Role        Role    `json:"role"`
	Track       string  `json:"track"`
	ContestID   uint    `json:"contest_id"`
	CommandID   uint    `json:"command_id"`
}

type Role string

const (
	Capitan     Role = "CAPITAN"
	Participant Role = "Participant"
)
