package user

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/api/user"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	userGroup := router.Group("/user")
	userGroup.Post("/register", user.Register)
	userGroup.Post("/login", user.Login)

	authorized := userGroup.Group("/authorized", middleware.UserAuth)
	authorized.Patch("/profile", user.ChangeUserProfile)
	authorized.Get("/profiles", user.ProfilesList)
	authorized.Get("/profiles/:id", user.ProfileByID)
	authorized.Post("/create_command", user.CreateCommand)
	authorized.Get("/contests/:profile_id", user.GetContestByID)
	authorized.Get("/commands", user.GetCommandByProfileID)
	authorized.Put("/commands", user.EditCommand)
	authorized.Get("/commands/leave", user.LeaveCommand)
	authorized.Put("/profiles", user.EditDescription)
	authorized.Get("/commands/matching", user.CommandMatching)

	authorized.Post("/command/createInvite", user.CreateInviteLink)
	authorized.Post("/command/useOInvite", user.UseInvite)
}
