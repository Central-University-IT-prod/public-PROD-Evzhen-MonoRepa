package entities

import "time"

type Invite struct {
	ID           uint      `json:"id"`
	CommandID    uint      `json:"command_id"`
	Code         string    `json:"code"`
	CreationDate time.Time `json:"creation_date"`
}
