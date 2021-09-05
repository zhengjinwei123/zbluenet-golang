package net

type NetMessage struct {
	MsgHead *MessageHead
	Info interface{}
}
