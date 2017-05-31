package main

import (
	"github.com/zhouhui8915/go-socket.io-client"
	"log"
	"time"
	"strconv"
)

var opts = &socketio_client.Options{
	//Transport:"polling",
	Transport:"websocket",
	Query:     make(map[string]string),
}
//opts.Query["user"] = "user"
//opts.Query["pwd"] = "pass"
var uri = "http://localhost:8000"

const (
	CLIENT_NUM=5000//同时多少客户端连接
	ROOM_NUM//每个房间最多多少人
)
func main() {
	clients:=make([]*socketio_client.Client,CLIENT_NUM)
	for i:=0;i<CLIENT_NUM;i++{
		clients[i],_=makeClients()
	}
	time.Sleep(time.Second*2)
	clients[0].Emit("printroom")
	for i:=0;i<CLIENT_NUM;i++{
		clients[i].Emit("talk","hello"+strconv.Itoa(1+i))
		time.Sleep(time.Microsecond*1)
	}

	time.Sleep(time.Minute*1)
	defer func(cs []*socketio_client.Client) {
		for _,client:=range cs{
			client.Emit("bye","disc")
		}
	}(clients)
}

func makeClients()(*socketio_client.Client,error){
	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return nil,err
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	client.On("connect", func() {
		log.Printf("on connect\n")
	})
	client.On("talk", func(msg string) {
		log.Printf("on talk:%v\n", msg)
	})
	client.On("bye", func(msg string) {
		log.Printf("on disconnect:%v\n",msg)
	})
	return client,nil

}