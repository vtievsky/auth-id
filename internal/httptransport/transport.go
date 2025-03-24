package httptransport

import (
	"github.com/vtievsky/auth-id/internal/conf"
	"github.com/vtievsky/auth-id/internal/services"
)

type Transport struct {
	conf     *conf.Config
	services *services.SvcLayer
}

func New(conf *conf.Config, svc *services.SvcLayer) *Transport {
	return &Transport{
		conf:     conf,
		services: svc,
	}
}
