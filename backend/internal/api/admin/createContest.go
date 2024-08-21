package admin

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/errorz"
	"github.com/gofiber/fiber/v2"
	"time"
)

type createContest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   int64  `json:"start_date"`
	EndDate     int64  `json:"end_date"`
	Field       string `json:"field"`
}

func validateCreateContest(data createContest) error {
	if !time.Unix(data.StartDate, 0).Before(time.Unix(data.EndDate, 0)) {
		return errorz.ErrorWrongTimeInterval
	}
	if !time.Now().Before(time.Unix(data.EndDate, 0)) {
		return errorz.ErrorWrongTimeInterval
	}
	return nil

}

func CreateContest(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	createContestData := new(createContest)
	if err = c.BodyParser(createContestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorz.CannotParseJSON.Error(),
		})
	}

	if err = validateCreateContest(*createContestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newContest := new(entities.Contest)
	newContest.Name = createContestData.Name
	newContest.Description = createContestData.Description
	newContest.StartDate = time.Unix(createContestData.StartDate, 0)
	newContest.EndDate = time.Unix(createContestData.EndDate, 0)
	newContest.Field = createContestData.Field
	newContest.AdminID = admin.ID

	database.DB.Create(newContest)

	return c.Status(fiber.StatusOK).JSON(newContest)
}
