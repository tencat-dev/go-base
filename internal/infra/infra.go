package infra

import (
	"github.com/goforj/wire"

	"github.com/tencat-dev/go-base/internal/infra/auth"
)

// ProviderSetInfra is infra providers.
var ProviderSetInfra = wire.NewSet(
	auth.NewJWTMaker,
)
