package main

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/nwehr/tictactoe/pkg/session"
)

var sessions []session.Session

func main() {

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
	log.Println("accepted new connection")
	defer conn.Close()

	for {
		msg := session.Message{}

		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
			log.Printf("could not decode message: %v\n", err)
			break
		}

		switch msg.Action {
		case session.Create:
			log.Println("new session")
			s := session.Session{
				ID:    msg.SessionID,
				State: msg.State,
				PlayerO: session.Player{
					ID:   "o",
					Conn: &conn,
				},
			}
			sessions = append(sessions, s)
			if err := gob.NewEncoder(conn).Encode(s.State); err != nil {
				log.Fatal("%v", err)
			}
		case session.Join:
			log.Println("join session")
			for i, s := range sessions {
				if s.ID == msg.SessionID {
					sessions[i].PlayerX = session.Player{
						ID:   "x",
						Conn: &conn,
					}

					for _, player := range []session.Player{sessions[i].PlayerO, sessions[i].PlayerX} {
						if err := gob.NewEncoder(*player.Conn).Encode(s.State); err != nil {
							log.Fatal("%v", err)
						}
					}

					break
				}

			}

		case session.Update:
			log.Println("update session")
			for _, s := range sessions {
				if s.ID == msg.SessionID {
					s.State = msg.State
					for _, player := range []session.Player{s.PlayerO, s.PlayerX} {
						log.Println("sending state to player %s", player.ID)
						if err := gob.NewEncoder(*player.Conn).Encode(s.State); err != nil {
							log.Fatal("%v", err)
						}
					}
				}
			}
		}
	}

}
