package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/nwehr/tictactoe/pkg/session"
	"github.com/nwehr/tictactoe/pkg/state"
)

func main() {
	host := prompt("Server host:port")
	action := prompt("Session create|join")
	sessionID := prompt("Session ID")
	player := "o"

	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal("could not connect to server: %v", err)
	}

	defer conn.Close()

	msg := session.Message{
		SessionID: sessionID,
		Action:    session.Create,
		State:     state.InitialGameState(),
	}

	if action == "join" {
		msg.Action = session.Join
		player = "x"
	}

	if err := gob.NewEncoder(conn).Encode(msg); err != nil {
		log.Fatal("could not encode/send msg: %v", err)
	}

	for {
		s := state.GameState{}
		if err := gob.NewDecoder(conn).Decode(&s); err != nil {
			log.Fatal("could not read/decode game state: %v", err)
		}

		s.Print()

		if p, won := s.Winner(); won == true {
			fmt.Printf("player %s wins!\n", p)
			break
		}

		if len(s.AvailableCaptures()) == 0 {
			fmt.Printf("tie!")
			break
		}

		if s.CurrentPlayer() == player {
			update := session.Message{
				SessionID: sessionID,
				Action:    session.Update,
				State:     s.Capture(prompt(fmt.Sprintf("%s capture", player))),
			}

			if err := gob.NewEncoder(conn).Encode(update); err != nil {
				log.Fatal("could not write/encode game state: %v", err)
			}
		}
	}

}

func prompt(question string) string {
	fmt.Printf("%s: ", question)

	response, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	return strings.TrimSpace(response)
}
