package config

import (
	"context"

	"github.com/cvartan/goconfig"
)

func GetConfig(ctx context.Context) *goconfig.Configuration {
	l := ctx.Value("config")
	if l == nil {
		panic("configuration is not defined in this context")
	}
	logger, ok := l.(*goconfig.Configuration)
	if !ok {
		panic("returned value from context for config is not *goconfig.Configuration")
	}
	return logger
}
