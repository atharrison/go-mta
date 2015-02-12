package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"strings"
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
		fmt.Println("Accepting new Connection. Starting new handleNewConnection goroutine.")
		go handleNewConnection(c)
	}
}

func handleNewConnection(conn net.Conn) {
	// handle an SMTP Conversation
//	msg := make([]byte, 1)
	// Connect

	// HELO
	fmt.Println("RESPONDING 250 HELO")
	conn.Write([]byte("250 HELO localhost\r\n"))
	status, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Received", status)

	// MAIL FROM
	conn.Write([]byte("250 OK\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Received", status)

	// RCPT TO
	conn.Write([]byte("250 OK\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Received", status)


	// DATA
	conn.Write([]byte("250 OK\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Received", status)

	conn.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))

	endOfData := false
	for {
		msg := make([]byte, 256)

		bytesRead, err := conn.Read(msg)
		data := string(msg[:])
		fmt.Println("Received", data)
		lines := strings.Split(data, "\r\n")

		fmt.Println("Last line: [", lines[len(lines)-1], "]")
		for _, line := range lines {
			if line == "." {
				endOfData = true
				break
			}
		}
		if endOfData {
			fmt.Println("Received terminating line, done reading DATA")
			break
		}
		if bytesRead == 0 {
			fmt.Println("No bytes read, finished reading DATA.")
			break
		}
		if err != nil {
			fmt.Println("Err: ", err)
			break
		}
	}
	fmt.Println("Acknowledging end of DATA.")

	// QUIT
	conn.Write([]byte("250 OK\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Received", status)

	conn.Write([]byte("221 localhost Service closing transmission channel.\r\n"))
	conn.Close()
}

func main() {
	fmt.Println("Starting the Go MTA Server.\n")

	go server()

	var input string
	fmt.Scanln(&input) // Die after input is read.

}
