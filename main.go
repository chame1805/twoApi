package main
import (
	"fmt"
	"log"
	"net/http"
	"reserva/src/reservas/Mensaje/infraestructure/controllers"
	"reserva/src/reservas/Mensaje/infraestructure/database"
	"reserva/src/reservas/Mensaje/infraestructure/routes"
	ws "reserva/src/reservas/Mensaje/infraestructure/socket"

	usecases "reserva/src/reservas/Mensaje/application/useCases"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Iniciando API...")

	hub := ws.NewHub()
	go hub.Run()

	rmq, err := database.NewRabbitMQ()
	if err != nil {
		log.Fatal("No se pudo conectar a RabbitMQ:", err)
	}
	defer rmq.Close()

	mensajeUseCase := usecases.NewProcesarMensajeUseCase(rmq)

	
	mensajeController := controllers.NewMensajeController(mensajeUseCase, hub)

	router := mux.NewRouter()

	router.HandleFunc("/mensaje", mensajeController.RecibirMensaje).Methods("POST")
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	routes.SetupRoutes(router, mensajeController)

	port := ":8081"
	fmt.Println("API escuchando en", port)
	log.Fatal(http.ListenAndServe(port, router))
}
