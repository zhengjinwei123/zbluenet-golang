package net

type TcpSession interface {
	OnInit(*TcpReactor)
	OnRecvMessage(*NetMessage)
	OnClose()
}