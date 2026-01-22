package service

import "github.com/goforj/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService)
