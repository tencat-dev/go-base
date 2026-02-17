package authz

import (
	"github.com/goforj/wire"
)

// ProviderSetAuthz is authz providers.
var ProviderSetAuthz = wire.NewSet(
	NewAuthzRegistry,
	NewAuthzMiddleware,
)
