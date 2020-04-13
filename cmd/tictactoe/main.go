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
	host := query("Server host:port")
	action := query("Session create|join")
	sessionID := query("Session ID")
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

		if s.AvailableMoves() == 0 {
			fmt.Printf("tie!")
			break
		}

		if s.NextTurn() == player {
			update := session.Message{
				SessionID: sessionID,
				Action:    session.Update,
				State:     s.Move(getMove(s)),
			}

			if err := gob.NewEncoder(conn).Encode(update); err != nil {
				log.Fatal("could not write/encode game state: %v", err)
			}
		}
	}

}

func play(s state.GameState) {
	s.Print()

	if player, won := s.Winner(); won == true {
		fmt.Printf("player %s wins!", player)
		os.Exit(0)
	}

	play(s.Move(getMove(s)))
}

func getMove(s state.GameState) string {
	fmt.Printf("%s move: ", s.NextTurn())
	move, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.ToUpper(strings.TrimSpace(move))
}

func query(question string) string {
	fmt.Printf("%s: ", question)

	response, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	return strings.TrimSpace(response)
}
