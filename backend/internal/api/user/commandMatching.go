package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"sort"
)

type CommandResponse struct {
	Type          string                `json:"type"`
	AveragePoints float64               `json:"average_points"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	CapitanTG     string                `json:"capitan_tg"`
	Participants  []ParticipantResponse `json:"participants"`
}

type ParticipantResponse struct {
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}

type CapitanResponse struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsCapitan   bool    `json:"is_capitan"`
	Rating      float64 `json:"rating"`
}

func CommandMatching(c *fiber.Ctx) error {
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
	var userProfiles []entities.Profile
	if err = database.DB.Where("user_login = ?", profile.UserLogin).Find(&userProfiles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}
	var avgProfile float64
	var sum, count float64
	for _, userProfile := range userProfiles {
		sum += userProfile.CurrPoints
		count++
	}
	avgProfile = sum / count

	switch profile.Role {
	case entities.Participant:
		var commands []entities.Command
		if err = database.DB.Where("contest_id = ? and approved = true", profile.ContestID).Find(&commands).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		commandsResponse := make([]CommandResponse, len(commands))

		for i, command := range commands {
			var participants []entities.Profile
			if err = database.DB.Where("command_id = ?", command.ID).Find(&participants).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}
			var capitan entities.Profile
			if err = database.DB.Where("id = ?", command.OwnerID).First(&capitan).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}

			commandsResponse[i].Participants = make([]ParticipantResponse, len(participants))
			commandsResponse[i].Type = string(entities.Participant)
			commandsResponse[i].Name = command.Name
			commandsResponse[i].Description = command.Description
			commandsResponse[i].CapitanTG = capitan.TG
			var avgSum, avgCount float64
			for j, participant := range participants {
				commandsResponse[i].Participants[j].Name = participant.Name
				var sum, count float64

				var userProfiles []entities.Profile
				if err = database.DB.Where("user_login = ?", participant.UserLogin).Find(&userProfiles).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status":  false,
						"message": err.Error(),
					})
				}
				for _, userProfile := range userProfiles {
					sum += userProfile.CurrPoints
					count++
				}
				sum += participant.PrevPoints
				count++
				commandsResponse[i].Participants[j].Rating = sum / count
				avgSum += sum / count
				avgCount++
			}
			commandsResponse[i].AveragePoints = avgSum / avgCount
		}

		sort.Slice(commandsResponse, func(i, j int) bool {
			return commandsResponse[i].AveragePoints-avgProfile < commandsResponse[j].AveragePoints-avgProfile
		})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"body": fiber.Map{
				"type":     entities.Participant,
				"commands": commandsResponse,
			},
		})
	case entities.Capitan:
		var participants []entities.Profile
		if err = database.DB.Where("contest_id = ? and command_id = 0", profile.ContestID).Find(&participants).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		capitanResponse := make([]CapitanResponse, len(participants))
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
			if participant.Role == entities.Capitan {
				isCapitan = true
			}

			capitanResponse[i] = CapitanResponse{
				Type:        string(entities.Capitan),
				Name:        participant.Name,
				Description: participant.Description,
				IsCapitan:   isCapitan,
				Rating:      sumProfile / countProfile,
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"body": fiber.Map{
				"type":         entities.Capitan,
				"participants": capitanResponse,
			},
		})
	}
	return nil
}
