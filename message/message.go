package message

import "github.com/pqnguyen/simple-chatapp/types"

type Message struct {
	UID    int              `json:"uid"`
	Action types.ActionType `json:"action"`
}

type Register struct {
	Message
}

func NewRegister(uid int) *Register {
	return &Register{Message{
		UID:    uid,
		Action: types.Register,
	}}
}

type Unregister struct {
	Message
}

func NewUnregister(uid int) *Unregister {
	return &Unregister{Message{
		UID:    uid,
		Action: types.Unregister,
	}}
}

type Talk struct {
	Message
	Content string `json:"content"`
	To      int    `json:"to"`
}

func NewTalk(uid int, to int, content string) *Talk {
	return &Talk{
		Message: Message{
			UID:    uid,
			Action: types.Talk,
		},
		To:      to,
		Content: content,
	}
}

type Ping struct {
	Content string `json:"content"`
}

func NewPing() *Ping {
	return &Ping{"Ping"}
}
