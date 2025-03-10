package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reserva/src/reservas/Mensaje/application/useCases"
	"reserva/src/reservas/Mensaje/domain/entity"
	"reserva/src/reservas/Mensaje/domain/repository"
	"reserva/src/reservas/Mensaje/infraestructure/socket"
)

// DTO que permite recibir tanto "message" como "contenido"
type mensajeDTO struct {
	Message   string `json:"message"`
	Contenido string `json:"contenido"`
}

type MensajeController struct {
	useCase *useCases.ProcesarMensajeUseCase
	wsHub   *socket.Socket // Corrección en el tipo de wsHub
}

// Constructor de MensajeController
func NewMensajeController(rabbitRepo repository.RabbitMQRepository, wsHub *socket.Socket) *MensajeController {
	return &MensajeController{
		useCase: useCases.NewProcesarMensajeUseCase(rabbitRepo),
		wsHub:   wsHub, // Inicialización correcta de wsHub
	}
}

// Método para recibir y procesar mensajes
func (c *MensajeController) RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("JSON recibido:", string(body))

	// Deserializamos en un DTO que admita tanto "message" como "contenido"
	var dto mensajeDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Mapear el DTO a la entidad de dominio
	// Si "contenido" está vacío, usamos "message"
	contenido := dto.Contenido
	if contenido == "" {
		contenido = dto.Message
	}

	// Asignar un ID por defecto o generar uno (aquí usamos "default-id")
	mensaje := entity.Mensaje{
		ID:        "default-id",
		Contenido: contenido,
	}

	// Imprimir el mensaje después de mapearlo
	fmt.Println("Mensaje después de Unmarshal y mapeo:", mensaje)

	// Ejecutar el use case para enviar el mensaje a RabbitMQ
	if err := c.useCase.Execute(mensaje); err != nil {
		http.Error(w, "Error al procesar el mensaje", http.StatusInternalServerError)
		return
	}

	// Convertir el mensaje a JSON para emitirlo vía WebSocket
	messageBytes, err := json.Marshal(mensaje)
	if err != nil {
		http.Error(w, "Error al serializar el mensaje", http.StatusInternalServerError)
		return
	}

	// Emitir el mensaje a todos los clientes conectados
	c.wsHub.Broadcast(messageBytes)
	fmt.Println("Mensaje procesado y emitido vía WebSocket")

	// Responder al cliente que el mensaje fue recibido y emitido
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mensaje recibido y emitido correctamente"))
}
