package middleware

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/util"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func UserAuth(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	if len(c.GetReqHeaders()["Authorization"]) > 0 {
		authHeader := c.GetReqHeaders()["Authorization"][0]
		if authHeader == "" {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "auth header is empty",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "invalid auth header",
			})
		}

		_, err := util.ParseJwt(parts[1], util.RoleUser)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "invalid auth header",
			})
		}
	} else {
		c.Status(fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "auth header is empty",
		})
	}

	return c.Next()
}

func AdminAuth(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	if len(c.GetReqHeaders()["Authorization"]) > 0 {
		authHeader := c.GetReqHeaders()["Authorization"][0]
		if authHeader == "" {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "auth header is empty",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "invalid auth header",
			})
		}

		_, err := util.ParseJwt(parts[1], util.RoleAdmin)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "invalid auth header",
			})
		}

		return c.Next()

	} else {
		c.Status(fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "auth header is empty",
		})
	}

}
