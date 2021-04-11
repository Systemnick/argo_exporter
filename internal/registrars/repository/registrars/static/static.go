package static

import (
	"context"
	"net"
	"time"

	"github.com/rs/zerolog"

	"github.com/Systemnick/argo_exporter/internal/registrars/repository/registrars"
)

type Repository struct {
	log *zerolog.Logger
}

func NewRepository(logger *zerolog.Logger) *Repository {
	return &Repository{
		log: logger,
	}
}

func (r *Repository) List(context.Context) (device []registrars.Device, err error) {
	timeStart := time.Now()
	r.log.Debug().Msg("Static registrar list")

	device = []registrars.Device{
		// {IP: net.IP{127, 0, 0, 1}},
		{IP: net.IP{10, 0, 0, 6}},
		{IP: net.IP{10, 0, 0, 10}},
		{IP: net.IP{10, 0, 0, 11}},
		{IP: net.IP{10, 0, 0, 12}},
		{IP: net.IP{10, 0, 0, 13}},
		{IP: net.IP{10, 0, 0, 44}},
		{IP: net.IP{10, 0, 0, 48}},
		{IP: net.IP{10, 0, 0, 71}},
		{IP: net.IP{10, 0, 0, 72}},
		{IP: net.IP{10, 0, 0, 73}},
		{IP: net.IP{10, 0, 0, 74}},
		{IP: net.IP{10, 0, 0, 75}},
		{IP: net.IP{10, 0, 0, 76}},
		{IP: net.IP{10, 0, 0, 77}},
		{IP: net.IP{10, 0, 0, 78}},
		{IP: net.IP{10, 0, 0, 79}},
		{IP: net.IP{10, 0, 0, 80}},
		{IP: net.IP{10, 0, 0, 85}},
	}

	r.log.Debug().Dur("duration", time.Since(timeStart)).Msg("Static registrar list finished")

	return device, nil
}
