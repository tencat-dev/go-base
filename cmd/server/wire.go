//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	conf "github.com/tencat-dev/go-base/codegen/proto"
	"github.com/tencat-dev/go-base/internal/biz"
	"github.com/tencat-dev/go-base/internal/data"
	"github.com/tencat-dev/go-base/internal/server"
	"github.com/tencat-dev/go-base/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goforj/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
