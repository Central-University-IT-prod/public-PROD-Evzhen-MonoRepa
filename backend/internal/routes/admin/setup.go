package admin

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/api/admin"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) {
	adminGroup := router.Group("/admin")
	adminGroup.Post("/login", admin.Login)
	adminGroup.Post("/reg", admin.Registration)

	authorizedAdminGroup := adminGroup.Group("/authorized", middleware.AdminAuth)
	authorizedAdminGroup.Get("/getAdmin", admin.GetAdmin)
	authorizedAdminGroup.Get("/getContests", admin.GetContests)
	authorizedAdminGroup.Get("/getCommands/:contest_id", admin.GetCommands)
	authorizedAdminGroup.Get("/getTeammates/:command_id", admin.GetCommandParticipants)

	authorizedAdminGroup.Post("/contests/create", admin.CreateContest)
	authorizedAdminGroup.Patch("/contest/change", admin.ChangeContest)
	authorizedAdminGroup.Post("/contests/loadProfiles", admin.LoadProfiles)
	authorizedAdminGroup.Patch("/contests/setTeamLimit", admin.SetTeamLimit)
	authorizedAdminGroup.Get("/contest/:contest_id/profiles", admin.GetProfiles)
	authorizedAdminGroup.Delete("/profile/:id", admin.DeleteProfile)

	authorizedAdminGroup.Patch("/contest/command/addProfile", admin.AddProfileToCommand)
	authorizedAdminGroup.Patch("/contest/command/removeProfile", admin.RemoveProfileFromCommand)
	authorizedAdminGroup.Delete("/contest/command", admin.DeleteCommand)
	authorizedAdminGroup.Patch("/contest/command/:command_id/approve/:approved", admin.ApproveCommand)

	authorizedAdminGroup.Post("/contest/end", admin.EndTheContest)
}
