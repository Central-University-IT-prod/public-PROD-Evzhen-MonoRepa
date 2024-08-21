package admin

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/errorz"
	"github.com/gofiber/fiber/v2"
)

type setTeamLimit struct {
	ContestID    uint `json:"contest_id"`
	MinTeammates uint `json:"min_teammates"`
	MaxTeammates uint `json:"max_teammates"`
}

func SetTeamLimit(c *fiber.Ctx) error {
	var (
		contest          entities.Contest
		setTeamLimitData setTeamLimit
	)
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	if err := c.BodyParser(&setTeamLimitData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorz.CannotParseJSON.Error(),
		})
	}

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err = database.DB.Where("id = ?", setTeamLimitData.ContestID).Find(&contest).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": errorz.ErrorContestNotFound.Error(),
		})
	}
	if contest.MinTeammates != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorz.ErrorTeamLimitAlreadyExists.Error(),
		})
	}
	if contest.MaxTeammates != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorz.ErrorTeamLimitAlreadyExists.Error(),
		})
	}

	if admin.ID != contest.AdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": errorz.ErrorMissingAccess.Error(),
		})
	}

	contest.MinTeammates = setTeamLimitData.MinTeammates
	contest.MaxTeammates = setTeamLimitData.MaxTeammates
	database.DB.Save(&contest)

	return c.Status(fiber.StatusOK).JSON(contest)
}
