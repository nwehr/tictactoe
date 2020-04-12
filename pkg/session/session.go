package session

import (
	"net"

	"github.com/nwehr/tictactoe/pkg/state"
)

type Session struct {
	ID    string
	State state.GameState

	PlayerX Player
	PlayerO Player
}

type Player struct {
	ID   string
	Conn *net.Conn
}
