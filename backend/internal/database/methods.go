package database

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/util"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetUserByAuthorizationToken(c *fiber.Ctx) (entities.User, error) {
	authHeader := c.GetReqHeaders()["Authorization"][0]

	id, err := util.ParseJwt(strings.Split(authHeader, " ")[1], util.RoleUser)
	if err != nil {
		return entities.User{}, err
	}

	var user entities.User
	if err = DB.Where("id = ?", id).Find(&user).Error; err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func GetAdminByAuthorizationToken(c *fiber.Ctx) (entities.Admin, error) {
	authHeader := c.GetReqHeaders()["Authorization"][0]

	id, err := util.ParseJwt(strings.Split(authHeader, " ")[1], util.RoleAdmin)
	if err != nil {
		return entities.Admin{}, err
	}

	var admin entities.Admin
	if err = DB.Where("id = ?", id).Find(&admin).Error; err != nil {
		return entities.Admin{}, err
	}

	return admin, nil
}

func GetCommandParticipants(commandId uint) []entities.Profile {
	var profiles []entities.Profile
	DB.Where("command_id = ?", commandId).Find(&profiles)
	return profiles
}
