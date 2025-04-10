package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}


func NewRabbitMQ() (*RabbitMQ, error) {
	
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del sistema")
	}

	
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	log.Println("Conectando a RabbitMQ en:", rabbitURL)


	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Println("Error al conectar con RabbitMQ:", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Error al abrir un canal RabbitMQ:", err)
		return nil, err
	}

	// Declarar la cola
	q, err := ch.QueueDeclare("mensajes", true, false, false, false, nil)
	if err != nil {
		log.Println("Error al declarar la cola:", err)
		return nil, err
	}

	return &RabbitMQ{conn: conn, channel: ch, queue: q}, nil
}


func (r *RabbitMQ) PublishMessage(message string) error {
	err := r.channel.Publish(
		"", r.queue.Name, false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Println("Error al enviar mensaje a RabbitMQ:", err)
		return err
	}

	fmt.Println("Mensaje enviado a RabbitMQ:", message)
	return nil
}

// Close cierra la conexión y el canal
func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
