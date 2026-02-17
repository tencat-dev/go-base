//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goforj/wire"

	"github.com/tencat-dev/go-base/internal/authz"
	"github.com/tencat-dev/go-base/internal/biz"
	"github.com/tencat-dev/go-base/internal/conf"
	"github.com/tencat-dev/go-base/internal/data"
	"github.com/tencat-dev/go-base/internal/infra"
	"github.com/tencat-dev/go-base/internal/server"
	"github.com/tencat-dev/go-base/internal/service"
)

// wireApp init kratos application.
func wireApp(context.Context, *conf.Bootstrap, log.Logger, *log.Helper) (*kratos.App, func(), error) {
	panic(wire.Build(
		ProviderSetConfig,
		data.ProviderSetData,
		authz.ProviderSetAuthz,
		biz.ProviderSetBiz,
		infra.ProviderSetInfra,
		service.ProviderSetService,
		server.ProviderSetServer,
		newApp,
	))
}
