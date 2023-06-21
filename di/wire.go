//go:build wireinject
// +build wireinject

package di

import (
	"blumer-ms-refers/handler"

	"github.com/google/wire"
)

func Initialize() (*handler.Handler, error) {
	wire.Build(stdSet)

	return &handler.Handler{}, nil
}
