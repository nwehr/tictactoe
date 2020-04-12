package session

import "github.com/nwehr/tictactoe/pkg/state"

type Message struct {
	SessionID string
	Action    Action
	State     state.GameState
}
