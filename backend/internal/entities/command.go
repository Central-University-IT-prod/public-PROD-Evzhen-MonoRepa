package entities

type Command struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//Participants []Profile `json:"participants" gorm:"foreignKey:CommandID; constraint:OnDelete:CASCADE"`
	ContestID uint `json:"contest_id"`
	OwnerID   uint `json:"owner_id"`
	Approved  bool `json:"approved"`
}
