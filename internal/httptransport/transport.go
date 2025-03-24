package httptransport

import (
	"github.com/vtievsky/auth-id/internal/conf"
	"github.com/vtievsky/auth-id/internal/services"
	"go.uber.org/zap"
)

type Transport struct {
	conf     *conf.Config
	logger   *zap.Logger
	services *services.SvcLayer
}

func New(conf *conf.Config, logger *zap.Logger, svc *services.SvcLayer) *Transport {
	return &Transport{
		conf:     conf,
		logger:   logger,
		services: svc,
	}
}
