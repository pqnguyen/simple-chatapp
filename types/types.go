package types

type ActionType int

const (
	Register ActionType = iota
	Unregister
	Talk
)

type MessageStatus int

const (
	Unread MessageStatus = iota
	Read
)
