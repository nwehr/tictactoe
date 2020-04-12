package session

type Action int

const (
	Create Action = iota
	Join
	Update
)
