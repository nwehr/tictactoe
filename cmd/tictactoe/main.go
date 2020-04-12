package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nwehr/tictactoe/pkg/state"
)

func main() {
	play(state.InitialGameState())
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
