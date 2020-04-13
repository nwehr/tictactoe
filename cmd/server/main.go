package main

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/nwehr/tictactoe/pkg/session"
)

var sessions map[string]*session.Session

func main() {
	sessions = map[string]*session.Session{}

	l, err := net.Listen("tcp", "0.0.0.0:3333")
	if err != nil {
		log.Fatal("could not setup socket: %v", err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("could not setup new connection: %v", err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer func() {
		log.Println("closing connection")
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}()

	for {
		msg := session.Message{}

		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
			log.Printf("could not decode message: %v\n", err)
			break
		}

		sess, ok := sessions[msg.SessionID]
		if !ok {
			sess = &session.Session{
				ID:    msg.SessionID,
				State: msg.State,
			}
			sessions[msg.SessionID] = sess
		}

		switch msg.Action {
		case session.Create:
			createSession(msg, &conn)
		case session.Join:
			joinSession(msg, &conn)
		case session.Update:
			updateSession(msg)
		}

		for _, player := range []session.Player{sess.PlayerO, sess.PlayerX} {
			if player.Conn == nil {
				log.Printf("player %s has a nil connection\n", player.ID)
				continue
			}

			if err := gob.NewEncoder(*player.Conn).Encode(sess.State); err != nil {
				log.Printf("could not write/encode state to player: %v", err)
			}
		}

	}

}

func createSession(msg session.Message, conn *net.Conn) {
	if sess, ok := sessions[msg.SessionID]; ok {
		sess.PlayerO = session.Player{
			ID:   "o",
			Conn: conn,
		}
	}
}

func joinSession(msg session.Message, conn *net.Conn) {
	if sess, ok := sessions[msg.SessionID]; ok {
		sess.PlayerX = session.Player{
			ID:   "x",
			Conn: conn,
		}
	}
}

func updateSession(msg session.Message) {
	if sess, ok := sessions[msg.SessionID]; ok {
		sess.State = msg.State
	}
}
