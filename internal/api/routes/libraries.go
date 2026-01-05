package routes

import (
	"bizarre-vpn-api/internal/api/handlers"
	"bizarre-vpn-api/internal/api/middlewares"
	"bizarre-vpn-api/internal/services"
)

func LibrariesRoutes(customRouter *CustomRouter) {
	customRouter.routerGroup.Use(
		middlewares.AuthRequired(customRouter.log, customRouter.cfg),
		middlewares.AdminRoleRequired(customRouter.log),
	)

	librariesService := services.NewLibrariesService(
		customRouter.log,
		customRouter.storage.BackendTypesStorage,
		customRouter.storage.ProtocolsStorage,
		customRouter.storage.LnkProtocolsBackendTypesStorage,
	)

	librariesHandler := handlers.LibrariesHandler{
		Log:              customRouter.log,
		LibrariesService: librariesService,
	}

	customRouter.routerGroup.GET(
		"/backend-types",
		librariesHandler.GetBackendTypesList,
	)

	customRouter.routerGroup.GET(
		"/protocols",
		librariesHandler.GetProtocolsList,
	)

	customRouter.routerGroup.GET(
		"/protocols-by-backend-type-id",
		librariesHandler.GetProtocolsListByBackendTypeId,
	)
}
