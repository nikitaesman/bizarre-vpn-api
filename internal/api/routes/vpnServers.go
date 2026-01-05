package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/middlewares"
	"bizarre-vpn-api/internal/services"
)

func VpnServersRoutes(customRouter *CustomRouter) {
	customRouter._router.Use(
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
	)

	vpnServersService := services.NewVpnServersService(
		customRouter.log,
		customRouter.storage.VpnServersStorage,
	)

	vpnServersHandler := handlers.VpnServersHandler{
		Log:               customRouter.log,
		VpnServersService: vpnServersService,
	}

	customRouter.routerGroup.GET(
		"",
		vpnServersHandler.GetExpandedList,
	)

	customRouter.routerGroup.GET(
		"/:id",
		vpnServersHandler.GetExpandedItem,
	)

	customRouter.routerGroup.POST(
		"",
		vpnServersHandler.CreateItem,
	)

	customRouter.routerGroup.DELETE(
		"/:id",
		vpnServersHandler.DeleteItem,
	)
}
