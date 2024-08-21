package admin

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type requestRemove struct {
	ProfileID uint `json:"profile_id"`
}

func RemoveProfileFromCommand(c *fiber.Ctx) error {
	var (
		requestRemoveStruct requestRemove
		profile             entities.Profile
		contest             entities.Contest
	)
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	if err := c.BodyParser(&requestRemoveStruct); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "unable to parse body",
		})
	}

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err = database.DB.Where("id = ?", requestRemoveStruct.ProfileID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if err = database.DB.Where("id = ?", profile.ContestID).First(&contest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if admin.ID != contest.AdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "missing access",
		})
	}

	var command entities.Command
	err = database.DB.Where("id = ?", profile.CommandID).Find(&command).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "command with this id not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	switch profile.Role {
	case entities.Capitan:
		var teammates []entities.Profile
		if err = database.DB.Where("command_id = ?", command.ID).Find(&teammates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		for _, teammate := range teammates {
			teammate.CommandID = 0
			if err = database.DB.Save(&teammate).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}
		}

		profile.Role = entities.Participant
		if err = database.DB.Save(&profile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err = database.DB.Delete(&command).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	case entities.Participant:
		profile.CommandID = 0

		if err = database.DB.Save(&profile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

type requestAdd struct {
	ContestID uint   `json:"contest_id"`
	UserLogin string `json:"user_login"`
	CommandID uint   `json:"command_id"`
}

func AddProfileToCommand(c *fiber.Ctx) error {
	var (
		requestAddStruct requestAdd
		profile          entities.Profile
		contest          entities.Contest
		command          entities.Command
	)
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	if err := c.BodyParser(&requestAddStruct); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "unable to parse body",
		})
	}

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err = database.DB.Where("contest_id = ? AND user_login = ?", requestAddStruct.ContestID, requestAddStruct.UserLogin).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if profile.CommandID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user already in command",
		})
	}

	if err = database.DB.Where("id = ?", requestAddStruct.CommandID).First(&command).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if err = database.DB.Where("id = ?", profile.ContestID).First(&contest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if admin.ID != contest.AdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "missing access",
		})
	}

	profile.CommandID = requestAddStruct.CommandID
	database.DB.Save(&profile)

	return c.Status(fiber.StatusOK).JSON(profile)
}

type requestDelete struct {
	CommandID uint `json:"command_id"`
}

func DeleteCommand(c *fiber.Ctx) error {
	var (
		requestDeleteStruct requestDelete
		teammates           []entities.Profile
		command             entities.Command
		contest             entities.Contest
	)
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	if err := c.BodyParser(&requestDeleteStruct); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "unable to parse body",
		})
	}

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err = database.DB.Where("id = ?", requestDeleteStruct.CommandID).First(&command).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if err = database.DB.Where("id = ?", command.ContestID).First(&contest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if admin.ID != contest.AdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "missing access",
		})
	}

	if err = database.DB.Where("command_id = ?", command.ID).Find(&teammates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for i := range teammates {
		teammates[i].CommandID = 0
		teammates[i].Role = entities.Participant
		if err = database.DB.Save(&teammates[i]).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
	}

	if err = database.DB.Delete(&command).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
