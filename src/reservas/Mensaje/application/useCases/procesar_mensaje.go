package useCases

import (
	"encoding/json"
	"reserva/src/reservas/Mensaje/domain/entity"
	"reserva/src/reservas/Mensaje/domain/repository"
)

type ProcesarMensajeUseCase struct {
	rabbitRepo repository.RabbitMQRepository
}

func NewProcesarMensajeUseCase(rabbitRepo repository.RabbitMQRepository) *ProcesarMensajeUseCase {
	return &ProcesarMensajeUseCase{
		rabbitRepo: rabbitRepo,
	}
}

func (uc *ProcesarMensajeUseCase) Execute(contenido string) error {
	mensaje := entity.Mensaje{
		ID:        "default-id",
		Contenido: contenido,
	}

	jsonMsg, err := json.Marshal(mensaje)
	if err != nil {
		return err
	}

	return uc.rabbitRepo.PublishMessage(string(jsonMsg))
}
