package main

import (
	"fmt"
	"log"
	"net/http"
	"reserva/src/reservas/Mensaje/infraestructure/controllers"
	"reserva/src/reservas/Mensaje/infraestructure/database"
	"reserva/src/reservas/Mensaje/infraestructure/routes"
	ws "reserva/src/reservas/Mensaje/infraestructure/socket"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Iniciando API...")

	// Iniciar el Hub de WebSocket
	hub := ws.NewHub()
	go hub.Run()

	// Conectar a RabbitMQ (implementación de RabbitMQRepository)
	rmq, err := database.NewRabbitMQ()
	if err != nil {
		log.Fatal("No se pudo conectar a RabbitMQ:", err)
	}
	defer rmq.Close()

	// Inyectar la implementación de RabbitMQ y el hub en el controlador.
	mensajeController := controllers.NewMensajeController(rmq, hub)

	// Configurar rutas usando mux.
	router := mux.NewRouter()

	// Aquí debes configurar las rutas que quieres manejar.
	// Puedes agregar rutas adicionales como se muestra a continuación.

	// Ruta para procesar mensajes
	router.HandleFunc("/mensaje", mensajeController.RecibirMensaje).Methods("POST")

	// Ruta para las conexiones WebSocket.
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	// Agregar más rutas a través de SetupRoutes si es necesario.
	routes.SetupRoutes(router, mensajeController)

	port := ":8081"
	fmt.Println("API escuchando en", port)
	log.Fatal(http.ListenAndServe(port, router))
}
