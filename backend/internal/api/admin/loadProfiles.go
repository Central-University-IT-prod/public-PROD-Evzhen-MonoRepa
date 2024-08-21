package admin

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

type LoadProfilesRequest struct {
	ContestID uint `json:"contest_id"`
	Profiles  []struct {
		Name       string  `json:"name"`
		Login      string  `json:"login"`
		TG         string  `json:"tg"`
		Track      string  `json:"track"`
		PrevPoints float64 `json:"prev_points"`
		MaxPoints  float64 `json:"max_points"`
	} `json:"profiles"`
}

type Err struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Err   string `json:"err"`
}

func LoadProfiles(c *fiber.Ctx) error {
	var loadProfileRequest LoadProfilesRequest
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	if err := c.BodyParser(&loadProfileRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to parse body",
		})
	}

	var contest entities.Contest
	if err := database.DB.Where("id = ?", loadProfileRequest.ContestID).First(&contest).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "contest not found",
		})
	}

	errs := make([]Err, 0)
	profiles := make([]entities.Profile, 0)
	for i, profile := range loadProfileRequest.Profiles {
		var user entities.User
		if err := database.DB.Where("login = ?", profile.Login).First(&user).Error; err != nil {
			errs = append(errs, Err{
				ID:    i + 1,
				Login: profile.Login,
				Err:   err.Error(),
			})
			continue
		}

		var p entities.Profile
		err := database.DB.Where("user_login = ? AND contest_id = ?", profile.Login, contest.ID).First(&p).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			errs = append(errs, Err{
				ID:    i + 1,
				Login: profile.Login,
				Err:   err.Error(),
			})
			continue
		}

		if err == nil {
			errs = append(errs, Err{
				ID:    i + 1,
				Login: profile.Login,
				Err:   "profile already exists",
			})
			continue
		}

		fio := strings.Split(profile.Name, " ")
		if len(fio) != 3 {
			errs = append(errs, Err{
				ID:    i + 1,
				Login: profile.Login,
				Err:   "invalid fio",
			})
			continue
		}

		prof := entities.Profile{
			Name:       fio[1],
			Surname:    fio[0],
			Patronymic: fio[2],
			TG:         profile.TG,
			Role:       entities.Participant,
			ContestID:  contest.ID,
			PrevPoints: profile.PrevPoints / profile.MaxPoints * 5,
			Track:      profile.Track,
			UserLogin:  profile.Login,
		}

		if err = database.DB.Create(&prof).Error; err != nil {
			errs = append(errs, Err{
				ID:    i + 1,
				Login: profile.Login,
				Err:   err.Error(),
			})
		}
		profiles = append(profiles, prof)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"created_profiles": profiles,
		"errors":           errs,
	})
}
