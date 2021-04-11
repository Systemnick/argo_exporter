package useCase

import (
	"context"

	"github.com/Systemnick/argo_exporter/internal/domain/flat"
)

type Indications interface {
	GetAll(ctx context.Context) Indications
	GetByFlat(flat flat.Flat)
}
