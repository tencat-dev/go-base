package server

import (
	"github.com/goforj/wire"
)

// ProviderSetServer is server providers.
var ProviderSetServer = wire.NewSet(
	NewGRPCServer,
	NewHTTPServer,
	NewPprofServer,
)
