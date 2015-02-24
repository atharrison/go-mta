package main

import (
	"net"
)

type ConnectionListener struct {
	smtpServerChan chan *SmtpServer
}

func NewConnectionListener(ch chan *SmtpServer) (cl ConnectionListener) {
	return ConnectionListener{ch}
}

func (cl ConnectionListener) start() {
	// listen on a port
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		Error.Println(err)
		return
	}

	for {
		// accept a connection
		conn, err := listener.Accept()
		if err != nil {
			Error.Println(err)
			continue
		}
		// handle the connection
		Info.Println("Accepting new Connection, placing on SmtpServer on Channel.")
		server := &SmtpServer{conn}
		cl.smtpServerChan <- server
	}
}
