package socket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub administra las conexiones WebSocket y distribuye mensajes a todos los clientes.
type Socket struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

// NewHub crea e inicializa un nuevo Hub.
func NewHub() *Socket {
	return &Socket{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run inicia el ciclo del hub, escuchando los canales de registro, baja y broadcast.
func (h *Socket) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			fmt.Println("Nuevo cliente conectado")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
				fmt.Println("Cliente desconectado")
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast envía un mensaje a todos los clientes conectados.
func (h *Socket) Broadcast(message []byte) {
	h.broadcast <- message
}