package main

import (
	"log"
	"net"

	"github.com/nwehr/tictactoe/pkg/state"
)

var sessions map[string]state.GameState

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:3333")
	if err != nil {
		log.Fatal("could not setup socket: %v", err)
	}

	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Fatal("could not setup new connection: %v", err)
	}
	_ = conn

}

func handleConnection() {

}
