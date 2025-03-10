package socket


import(
	"net/http"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	// Permitir conexiones de cualquier origen (útil en desarrollo)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(h *Socket, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.register <- conn
}