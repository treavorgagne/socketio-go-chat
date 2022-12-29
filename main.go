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

	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")

		server.OnEvent("/", "new user", func(s socketio.Conn, name string, room string) {
			server.BroadcastToRoom("/", room, "chat sub msg", name + " entered the chat")
			s.Join(room)
		})

		// sends message to everyone in the chat
		server.OnEvent("/", "send message", func(s socketio.Conn, msg string, name string, room string) {
			server.BroadcastToRoom("/", room, "receive message", name +": "+ msg)
		});

		server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		})

		return nil
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	
	log.Println("Serving at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}