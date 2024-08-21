package admin

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type endContestRequest struct {
	ContestID uint `json:"contest_id"`
	Profiles  []struct {
		Login       string  `json:"login"`
		FinalPoints float64 `json:"final_points"`
		MaxPoints   float64 `json:"max_final_points"`
	} `json:"profiles"`
}

type errEndContest struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Err   string `json:"err"`
}

func EndTheContest(c *fiber.Ctx) error {
	var endContestData endContestRequest
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	if err := c.BodyParser(&endContestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to parse body",
		})
	}

	var contest entities.Contest
	if err := database.DB.Where("id = ?", endContestData.ContestID).First(&contest).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "contest not found",
		})
	}

	errs := make([]errEndContest, 0)
	profiles := make([]entities.Profile, 0)
	for i, profile := range endContestData.Profiles {
		var user entities.User
		if err := database.DB.Where("login = ?", profile.Login).First(&user).Error; err != nil {
			errs = append(errs, errEndContest{
				ID:    i + 1,
				Login: profile.Login,
				Err:   err.Error(),
			})
			continue
		}

		var p entities.Profile
		err := database.DB.Where("user_login = ? AND contest_id = ?", profile.Login, contest.ID).First(&p).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			errs = append(errs, errEndContest{
				ID:    i + 1,
				Login: profile.Login,
				Err:   err.Error(),
			})
			continue
		}

		if err != nil {
			errs = append(errs, errEndContest{
				ID:    i + 1,
				Login: profile.Login,
				Err:   "profile not found",
			})
			continue
		}

		p.CurrPoints = profile.FinalPoints / profile.MaxPoints * 5

		if err = database.DB.Where("user_login = ? AND contest_id = ?", profile.Login, contest.ID).Save(&p).Error; err != nil {
			errs = append(errs, errEndContest{
				ID:    i + 1,
				Login: profile.Login,
				Err:   err.Error(),
			})
			continue
		}
		profiles = append(profiles, p)
	}

	contest.End = true
	database.DB.Save(&contest)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"changed_profiles": profiles,
		"errors":           errs,
	})
}
