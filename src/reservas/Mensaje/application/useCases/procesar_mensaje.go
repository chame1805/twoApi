package useCases

import (
	"encoding/json"
	"fmt"
	"reserva/src/reservas/Mensaje/domain/entity"
	"reserva/src/reservas/Mensaje/domain/repository"
)

type ProcesarMensajeUseCase struct {
	rabbitRepo repository.RabbitMQRepository
}

func NewProcesarMensajeUseCase(rabbitRepo repository.RabbitMQRepository) *ProcesarMensajeUseCase {
	return &ProcesarMensajeUseCase{rabbitRepo: rabbitRepo}
}

func (pm *ProcesarMensajeUseCase) Execute(mensaje entity.Mensaje) error {
	fmt.Printf("Mensaje recibido en useCase: %+v\n", mensaje)

	if mensaje.ID == "" || mensaje.Contenido == "" {
		return fmt.Errorf("Mensaje inválido: ID o Contenido vacío")
	}

	// Convertir el mensaje a JSON
	mensajeJSON, err := json.Marshal(mensaje)
	if err != nil {
		return fmt.Errorf("Error al convertir mensaje a JSON: %s", err)
	}

	// Publicar el mensaje en RabbitMQ
	err = pm.rabbitRepo.PublishMessage(string(mensajeJSON))
	if err != nil {
		return fmt.Errorf("Error al enviar mensaje a la cola: %s", err)
	}

	// Publicar una notificación para que el servidor de polling sepa que hay un nuevo mensaje
	notificacion := "Nuevo mensaje disponible"
	err = pm.rabbitRepo.PublishMessage(notificacion)  // Esto es la notificación
	if err != nil {
		return fmt.Errorf("Error al enviar notificación a la cola: %s", err)
	}

	fmt.Println("Mensaje enviado a RabbitMQ correctamente:", string(mensajeJSON))
	fmt.Println("Notificación enviada a la cola de notificaciones")

	return nil
}


