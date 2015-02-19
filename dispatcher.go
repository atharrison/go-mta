package main

import (
	"net"
)

func handleEnvelope(env envelope) {
	Info.Println("Dispatcher dispatching message...")
	Debug.Println("Received Envelope:", env)

	// TODO: Determine where to send based on RcptTo domain
	conn, err := createNewSmtpConnection()
	if err != nil {
		Error.Println("Could not connect, Failed to send message.", err)
		// TODO Don't just drop the message on the floor.
	} else {
		smtpClientChan <- &SmtpClient{conn, env}
	}
}

func createNewSmtpConnection() (conn net.Conn, err error) {

	Info.Println("Creating new outbound SMTP connection to localhost:2525...")
	conn, err = net.Dial("tcp", "localhost:2525")
	return conn, err
}


