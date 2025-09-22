//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"ragx/app/internal/biz"
	"ragx/app/internal/conf"
	"ragx/app/internal/data"
	"ragx/app/internal/server"
	"ragx/app/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
