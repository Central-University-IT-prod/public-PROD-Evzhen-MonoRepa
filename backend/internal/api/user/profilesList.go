package user

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProfileResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  string  `json:"patronymic"`
	Description string  `json:"description"`
	PrevPoints  float64 `json:"prev_points"`
	CurrPoints  float64 `json:"curr_points"`
	Role        string  `json:"role"`
	ContestID   uint    `json:"contest_id"`
	Login       string  `json:"login"`
}

func ProfilesList(c *fiber.Ctx) error {
	var limit, offset int

	limit, _ = strconv.Atoi(c.Params("limit"))
	offset, _ = strconv.Atoi(c.Params("offset"))
	if limit == 0 {
		limit = 5
	}

	user, err := database.GetUserByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var profiles []entities.Profile
	if err = database.DB.Where("user_login = ?", user.Login).Limit(limit).Offset(offset).Find(&profiles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if len(profiles) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"body": fiber.Map{
				"profiles": []ProfileResponse{},
			},
		})
	}

	response := make([]ProfileResponse, len(profiles))
	for i, profile := range profiles {
		response[i] = ProfileResponse{
			ID:          profile.ID,
			Name:        profile.Name,
			Surname:     profile.Surname,
			Patronymic:  profile.Patronymic,
			Description: profile.Description,
			PrevPoints:  profile.PrevPoints,
			CurrPoints:  profile.CurrPoints,
			Role:        string(profile.Role),
			ContestID:   profile.ContestID,
			Login:       profile.UserLogin}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body": fiber.Map{
			"profiles": response,
		},
	})
}
