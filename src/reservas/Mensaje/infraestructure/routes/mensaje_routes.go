package routes

import (
	"reserva/src/reservas/Mensaje/infraestructure/controllers"

	"github.com/gorilla/mux"
	
)

func SetupRoutes(router *mux.Router, mensajeController *controllers.MensajeController) {
	router.HandleFunc("/citas", mensajeController.RecibirMensaje).Methods("POST")
}
