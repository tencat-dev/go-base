package biz

import "github.com/goforj/wire"

// ProviderSetBiz is biz providers.
var ProviderSetBiz = wire.NewSet(
	NewUserBiz,
	NewAuthBiz,
	NewAuthzBiz,
)
