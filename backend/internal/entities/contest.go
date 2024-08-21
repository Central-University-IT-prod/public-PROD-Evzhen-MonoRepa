package entities

import "time"

type Contest struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Field        string    `json:"field"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	MinTeammates uint      `json:"min_teammates"`
	MaxTeammates uint      `json:"max_teammates"`
	End          bool      `json:"end"`
	//Profiles    []Profile `json:"profiles" gorm:"foreignKey:ContestID constraint:OnDelete:CASCADE"`
	//Commands    []Command `json:"commands" gorm:"foreignKey:ContestID; constraint:OnDelete:CASCADE"`
	AdminID uint `json:"admin_id"`
}
