package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {

	fmt.Println("hello world")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./views/index.html")
    })

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./views/room.html")
    })

	var room string
	var username string

	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	server.OnEvent("/", "new user", func(s socketio.Conn, name string, newroom string) {
		room = newroom
		username = name
		s.Join(room)
		server.BroadcastToRoom("/", room, "chat sub msg", username + " entered the chat")
	})

	// sends message to everyone in the chat
	server.OnEvent("/", "send message", func(s socketio.Conn, msg string) {
		server.BroadcastToRoom("/", room, "receive message", username +": "+ msg)
	});

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		server.BroadcastToRoom("/", room, "chat sub msg", username + " left the chat")
		fmt.Println("closed:", reason, room, username)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	
	log.Println("Serving at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}