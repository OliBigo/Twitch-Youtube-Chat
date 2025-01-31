package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	upgrader         = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients          = make(map[*Client]bool)
	broadcast        = make(chan Message, 100)
	connectedYoutube string
	connectedTwitch  string
	twitchClient     *twitch.Client
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/youtube", handleChannelYoutube)
	mux.HandleFunc("/twitch", handleChannelTwitch)
	mux.HandleFunc("/ws", handleConnections)

	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)

	server := &http.Server{
		Addr:    port,
		Handler: enableCORS(mux),
	}

	go handleMessages()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    defer conn.Close()

    client := &Client{conn: conn, id: conn.RemoteAddr().String()}
    clients[client] = true

    log.Printf("New client connected: %s", client.id)

    for {
        var msg Message
        if err := conn.ReadJSON(&msg); err != nil {
            log.Printf("Client disconnected: %s", client.id)
            delete(clients, client)
            break
        }
        broadcast <- msg
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            if err := client.conn.WriteJSON(msg); err != nil {
                log.Printf("Error writing to client: %v", err)
                client.conn.Close()
                delete(clients, client)
            }
        }
    }
}