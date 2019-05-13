package types

type ActionType int

const (
	Register ActionType = iota
	Unregister
	Talk
)
