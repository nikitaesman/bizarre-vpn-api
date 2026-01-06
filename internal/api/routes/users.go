package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/middlewares"
	"bizarre-vpn-api/internal/services"
)

func UserRoutes(customRouter *CustomRouter) {
	userService := services.NewUserService(
		customRouter.log,
		customRouter.storage.UserStorage,
	)

	inviteLinksService := services.NewInviteLinksService(
		customRouter.log,
		customRouter.storage.InviteLinksStorage,
		customRouter.storage.LnkUserProviderStorage,
	)

	userHandler := handlers.UserHandler{
		Log:                customRouter.log,
		BotSharedData:      customRouter.botSharedData,
		UserService:        userService,
		InviteLinksService: inviteLinksService,
	}

	customRouter.AddGroup("/auth", AuthRoutes)

	customRouter.routerGroup.GET("/self",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		userHandler.GetUserDataHandler,
	)

	customRouter.routerGroup.GET("/",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		userHandler.GetUsersListHandler,
	)

	customRouter.routerGroup.POST("/",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		userHandler.CreateUser,
	)

	customRouter.routerGroup.PUT("/:id",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		userHandler.UpdateUser,
	)

	customRouter.routerGroup.DELETE("/:id",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		userHandler.DeleteUser,
	)

	customRouter.routerGroup.GET("/:id/invite-link",
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
		userHandler.CreateUserInviteLink,
	)
}
