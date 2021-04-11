package indications

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/Systemnick/argo_exporter/internal/domain/flat"
	"github.com/Systemnick/argo_exporter/internal/domain/registrar"
	"github.com/Systemnick/argo_exporter/internal/indications/repository/indications"
	"github.com/Systemnick/argo_exporter/internal/registrars/useCase"
)

type UseCase struct {
	registrars        useCase.Registrar
	indicationsFabric indications.Repository

	log *zerolog.Logger
}

func NewUseCase(logger *zerolog.Logger, registrarsUC useCase.Registrar, indicationFabric indications.Repository) *UseCase {
	return &UseCase{
		registrars:        registrarsUC,
		indicationsFabric: indicationFabric,

		log: logger,
	}
}

func (u *UseCase) Scrape(ctx context.Context, in *indications.MURMetrics) {
	timeStart := time.Now()
	u.log.Debug().Msg("Indications scrape")

	for _, reg := range u.registrars.List(ctx) {
		go func(reg *registrar.Registrar) {
			timeStart := time.Now()
			u.log.Debug().Msg("Indications scrape goroutine")

			if reg.IsPollingNow() {
				fmt.Printf("%s is still polling, skip...\n", reg.IP)
				return
			}

			err := u.indicationsFabric.Poll(reg.IP, in)
			if err != nil {
				err = fmt.Errorf("%s poll error: %w\n", reg.IP, err)
				fmt.Printf("%v\n", err)
			}

			u.log.Debug().Dur("duration", time.Since(timeStart)).Msg("Indications scrape goroutine finished")
		}(reg)
	}

	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// u.indicationsFabric.OpenConnections(func() {
	// 	wg.Done()
	// })
	// wg.Wait()
	//
	// wg.Add(1)
	// u.indicationsFabric.Poll(in, func() {
	// 	wg.Done()
	// })
	// wg.Wait()
	//
	// wg.Add(1)
	// u.indicationsFabric.CloseConnections(func() {
	// 	wg.Done()
	// })
	// wg.Wait()

	u.log.Debug().Dur("duration", time.Since(timeStart)).Msg("Indications scrape running")
}

func (u *UseCase) ScrapePeriodic(ctx context.Context, interval time.Duration) (*indications.MURMetrics, func()) {
	var in = indications.NewMURMetrics()

	ticker := time.NewTicker(interval)
	quit := make(chan struct{})

	cancelFunc := func() {
		quit <- struct{}{}
	}

	go func() {
		// First run
		go u.Scrape(ctx, in)

		for {
			select {
			case <-ticker.C:
				// fmt.Printf("\nTick...\n\n")
				go u.Scrape(ctx, in)
			case <-quit:
				ticker.Stop()

				return
			}
		}
	}()

	return in, cancelFunc
}

func (u *UseCase) GetByFlat(flat flat.Flat) {
	panic("implement me")
}
