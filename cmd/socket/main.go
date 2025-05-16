package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dwnwp/api-email/socket"
	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/ws", handleWS)
	http.HandleFunc("/notify", notifyUser)

	fmt.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	var upgrader websocket.Upgrader
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	socket.SetConn(conn)

	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.Println("Received:", string(msg))
	}
}

func notifyUser(w http.ResponseWriter, r *http.Request) {
	conn := socket.GetConn()

	if conn == nil {
		http.Error(w, "No connection found", http.StatusNotFound)
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte("Job done!"))
	if err != nil {
		log.Println("Send failed:", err)
		http.Error(w, "Send failed", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Message sent to client")
}
