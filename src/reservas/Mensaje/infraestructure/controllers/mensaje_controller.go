package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reserva/src/reservas/Mensaje/application/useCases"
	"reserva/src/reservas/Mensaje/infraestructure/socket"
)

type mensajeDTO struct {
	Message   string `json:"message"`
	Contenido string `json:"contenido"`
}

type MensajeController struct {
	useCase *useCases.ProcesarMensajeUseCase
	wsHub   *socket.Socket
}

func NewMensajeController(useCase *useCases.ProcesarMensajeUseCase, wsHub *socket.Socket) *MensajeController {
	return &MensajeController{
		useCase: useCase,
		wsHub:   wsHub,
	}
}

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

	var dto mensajeDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if dto.Contenido == "" {
		dto.Contenido = dto.Message
	}

	// Llamar al use case con el contenido
	if err := c.useCase.Execute(dto.Contenido); err != nil {
		http.Error(w, "Error al procesar el mensaje", http.StatusInternalServerError)
		return
	}

	// Emitir a WebSocket
	msgJSON, _ := json.Marshal(map[string]string{"contenido": dto.Contenido})
	c.wsHub.Broadcast(msgJSON)

	fmt.Println("Mensaje procesado y emitido vía WebSocket")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mensaje recibido y emitido correctamente"))
}
