package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GetCommandByProfileIDResponse struct {
	ID            uint          `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Participants  []Participant `json:"participants"`
	ContestID     uint          `json:"contest_id"`
	OwnerID       uint          `json:"owner_id"`
	AveragePoints float64       `json:"average_points"`
}

type Participant struct {
	IsCapitan bool    `json:"is_capitan"`
	UserLogin string  `json:"user_login"`
	Rating    float64 `json:"rating"`
}

func GetCommandByProfileID(c *fiber.Ctx) error {
	profileID := c.Query("profile_id")

	var profile entities.Profile
	err := database.DB.Where("id = ?", profileID).First(&profile).Error
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

	if profile.CommandID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": true,
			"body": fiber.Map{
				"no_command": true,
			},
		})
	}

	var command entities.Command
	err = database.DB.Where("id = ?", profile.CommandID).First(&command).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "command with this id not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var participants []entities.Profile
	if err = database.DB.Where("command_id = ?", command.ID).Find(&participants).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}
	response := GetCommandByProfileIDResponse{
		ID:           command.ID,
		Name:         command.Name,
		Description:  command.Description,
		Participants: make([]Participant, len(participants)),
		ContestID:    command.ContestID,
		OwnerID:      command.OwnerID,
	}

	var sum, count float64
	for i, participant := range participants {
		var profiles []entities.Profile
		if err = database.DB.Where("user_login = ?", participant.UserLogin).First(&profiles).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		var sumProfile, countProfile float64
		for _, uProfile := range profiles {
			sumProfile += uProfile.CurrPoints
			countProfile++
		}

		var isCapitan bool
		if profile.Role == entities.Capitan {
			isCapitan = true
		}

		response.Participants[i] = Participant{
			IsCapitan: isCapitan,
			UserLogin: participant.UserLogin,
			Rating:    sumProfile / countProfile,
		}
		sum += participant.CurrPoints
		count++
	}
	response.AveragePoints = sum / count

	var isCapitan bool
	if profile.Role == entities.Capitan {
		isCapitan = true
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body": fiber.Map{
			"no_command": false,
			"is_capitan": isCapitan,
			"command":    response,
		},
	})
}
