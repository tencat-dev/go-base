package service

import "github.com/goforj/wire"

// ProviderSetService is service providers.
var ProviderSetService = wire.NewSet(
	NewUserService,
	NewAuthService,
	NewAuthzService,
)
