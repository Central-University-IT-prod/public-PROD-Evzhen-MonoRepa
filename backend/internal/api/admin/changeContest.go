package admin

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/errorz"
	"github.com/gofiber/fiber/v2"
	"time"
)

type changeContestStruct struct {
	ContestID   uint   `json:"contest_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   int64  `json:"start_date"`
	EndDate     int64  `json:"end_date"`
}

func ChangeContest(c *fiber.Ctx) error {
	var (
		changeContestData changeContestStruct
		contest           entities.Contest
	)

	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err = c.BodyParser(&changeContestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorz.CannotParseJSON.Error(),
		})
	}

	if err = database.DB.Where("id = ?", changeContestData.ContestID).First(&contest).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	if admin.ID != contest.AdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "missing access",
		})
	}

	if time.Unix(changeContestData.StartDate, 0).After(time.Unix(changeContestData.EndDate, 0)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	if !time.Now().Before(time.Unix(changeContestData.EndDate, 0)) {
		contest.End = true
	} else {
		contest.End = false
	}

	contest.Name = changeContestData.Name
	contest.Description = changeContestData.Description
	contest.StartDate = time.Unix(changeContestData.StartDate, 0)
	contest.EndDate = time.Unix(changeContestData.EndDate, 0)

	database.DB.Save(&contest)

	return c.Status(fiber.StatusOK).JSON(contest)
}
