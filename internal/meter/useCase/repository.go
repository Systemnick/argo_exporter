package useCase

import "github.com/google/uuid"

type PollRepoInterface interface {
	GetCurrentIndication(registrarId uuid.UUID, meterId int16) ([]MeterInfo, error)
}
