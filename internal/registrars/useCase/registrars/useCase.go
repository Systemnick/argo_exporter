package registrars

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/Systemnick/argo_exporter/internal/domain/flat"
	"github.com/Systemnick/argo_exporter/internal/domain/registrar"
	"github.com/Systemnick/argo_exporter/internal/registrars/repository/registrars"
)

type UseCase struct {
	Registrar registrars.Repository

	log *zerolog.Logger
}

func NewUseCase(logger *zerolog.Logger, RegistrarsRepo registrars.Repository) *UseCase {
	return &UseCase{
		Registrar: RegistrarsRepo,

		log: logger,
	}
}

func (u *UseCase) List(ctx context.Context) []*registrar.Registrar {
	timeStart := time.Now()
	u.log.Debug().Msg("Registrar list")

	devices, err := u.Registrar.List(ctx)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	ips := make([]*registrar.Registrar, len(devices))
	for i, device := range devices {
		ips[i] = &registrar.Registrar{
			IP:         device.IP,
			Vendor:     "Argo",
			Model:      "MUR",
		}
	}

	u.log.Debug().
		Dur("duration", time.Since(timeStart)).
		Int("count", len(ips)).
		Msg("Registrar list finished")

	return ips
}

func (u *UseCase) GetByFlat(flat flat.Flat) {
	panic("implement me")
}
