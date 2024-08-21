package user

import (
	"errors"
	"fmt"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func CreateInviteLink(c *fiber.Ctx) error {
	var (
		profile entities.Profile
		command entities.Command
	)
	profileID := c.Query("profile_id")

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
			"status":  false,
			"message": "this profile does not belong to any of the commands",
		})
	}

	err = database.DB.Where("id = ?", profileID).First(&command).Error
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
	if profile.ID != command.OwnerID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": "missing access",
		})
	}

	newInvite := entities.Invite{
		CommandID:    command.ID,
		Code:         uuid.New().String(),
		CreationDate: time.Now(),
	}
	database.DB.Create(&newInvite)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
		"body": fiber.Map{
			"inviteLink": fmt.Sprintf("https://maxsiomin.dev/teamder/invite/%s", newInvite.Code),
		},
	})
}

func UseInvite(c *fiber.Ctx) error {
	var (
		invite       entities.Invite
		command      entities.Command
		participants []entities.Profile
		contest      entities.Contest
		profile      entities.Profile
	)

	inviteCode := c.Query("code")
	profileID := c.Query("profile_id")

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

	if profile.CommandID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "user already in command",
		})
	}

	err = database.DB.Where("id = ?", inviteCode).First(&invite).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "invite with this code not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if time.Since(invite.CreationDate) >= 24*time.Hour {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invite link expired",
		})
	}

	err = database.DB.Where("id = ?", invite.CommandID).First(&command).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "command not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if err = database.DB.Where("command_id = ?", command.ID).Find(&participants).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	err = database.DB.Where("id = ?", command.ContestID).First(&contest).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "contest not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if len(participants) >= int(contest.MaxTeammates) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "command already full",
		})
	}

	userInThisCommand := false
	for _, participant := range participants {
		if profile.ID == participant.ID {
			userInThisCommand = true
		}
	}

	if userInThisCommand {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "this user already in command",
		})
	}

	profile.CommandID = invite.CommandID
	database.DB.Save(&profile)
	database.DB.Delete(&invite)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
	})
}
