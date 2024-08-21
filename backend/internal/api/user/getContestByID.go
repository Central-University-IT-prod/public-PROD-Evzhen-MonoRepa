package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GetContestByIDResponse struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Field             string `json:"field"`
	ParticipantsCount uint   `json:"participants_count"`
	StartDate         int64  `json:"start_date"`
	EndDate           int64  `json:"end_date"`
	MinTeammates      uint   `json:"min_teammates"`
	MaxTeammates      uint   `json:"max_teammates"`
	AdminID           uint   `json:"admin_id"`
}

func GetContestByID(c *fiber.Ctx) error {
	id := c.Params("profile_id")

	var profile entities.Profile
	err := database.DB.Where("id = ?", id).First(&profile).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "profile with this id not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var contest entities.Contest
	err = database.DB.Where("id = ?", profile.ContestID).First(&contest).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "contest with this id not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var profiles []entities.Profile
	if err = database.DB.Where("contest_id = ?", contest.ID).Find(&profiles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	response := GetContestByIDResponse{
		ID:                contest.ID,
		Name:              contest.Name,
		Description:       contest.Description,
		Field:             contest.Field,
		ParticipantsCount: uint(len(profiles)),
		StartDate:         contest.StartDate.Unix(),
		EndDate:           contest.EndDate.Unix(),
		MinTeammates:      contest.MinTeammates,
		MaxTeammates:      contest.MaxTeammates,
		AdminID:           contest.AdminID,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body": fiber.Map{
			"contest": response,
		},
	})
}
