package useCase

import (
	"context"

	"github.com/Systemnick/argo_exporter/internal/domain/registrar"
)

type Registrar interface {
	List(ctx context.Context) []*registrar.Registrar
}
