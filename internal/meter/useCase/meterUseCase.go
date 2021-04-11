package useCase

import (
	"github.com/google/uuid"
)

type MeterUseCase struct {
	PollRepo PollRepoInterface
}

func NewMeterUseCase() *MeterUseCase {
	return &MeterUseCase{}
}

func (u MeterUseCase) GetCurrentIndications(registrarId uuid.UUID, meterId int16) FlatIndication {
	i, err := u.PollRepo.GetCurrentIndication(registrarId, meterId)
	if err != nil {
		return FlatIndication{}
	}

	return FlatIndication{
		Indications: i,
	}
}
