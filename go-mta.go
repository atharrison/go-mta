package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

func server() {
	// listen on a port
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// accept a connection
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// handle the connection
		go handleServerConnection(c)
	}
}

func handleServerConnection(c net.Conn) {
	// handle an SMTP Conversation
	msg := make([]byte, 1)

	// Connect
	for {
		_, _ = c.Read(msg)
		fmt.Println("Received", msg)
	}

	// HELO
	c.Write([]byte("HELO world"))

//	// MAIL FROM
//	_, err = c.Read(msg)
//	fmt.Println("Received", msg)
//
//	c.Write([]byte("250 OK"))
//
//	// RCPT TO
//	_, err = c.Read(msg)
//	fmt.Println("Received", msg)
//
//	c.Write([]byte("250 OK"))
//
//	// DATA
//	_, err = c.Read(msg)
//	fmt.Println("Received", msg)
//
//	c.Write([]byte("354 End data with <CR><LF>.<CR><LF>"))
//
//
//	c.Read(msg)

	c.Close()
}

func client() {
	// connect to the server
	c, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	// send the message
	msg := "Hello World"
	fmt.Println("Sending", msg)
	err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}

	c.Close()
}

func main() {
	fmt.Println("Starting the Go MTA Server.\n")

	go server()
//	go client()

	var input string
	fmt.Scanln(&input) // Die after input is read.

}
