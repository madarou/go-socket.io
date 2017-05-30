package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
	"unicode/utf8"
	"strconv"
)

var(
	room_num =0
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("ctx"+s.ID())
		room:=findRoom(s.ID())
		fmt.Println("connected:", s.ID(), " joined: "+ room)
		server.Rooms().Join(room,s)
		return nil
	})

	server.OnEvent("/", "talk", func(s socketio.Conn, msg string) {
		room:=findRoom(s.ID())
		fmt.Println(s.ID()+" talk:", msg, " in ", room)
		//s.Emit("reply", "reply "+msg)
		server.Rooms().BroadcastTo(s,room,"talk",msg)
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		//找到s所属的room
		room,_:=server.Rooms().Belong(s)
		server.Rooms().Leave(room,s)
		s.Close()
		return last
	})

	server.OnEvent("/","printroom", func(s socketio.Conn) {
		server.Rooms().List()
	})

	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	//http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func findRoom(id string)string{
	if utf8.RuneCountInString(id)==1{
		return "room"+strconv.Itoa(0)
	}
	if utf8.RuneCountInString(id)==2{
		return "room"+id[0:1]
	}
	if utf8.RuneCountInString(id)==3{
		return "room"+id[0:2]
	}
	return "room"
}