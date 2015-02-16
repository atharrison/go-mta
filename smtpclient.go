// Handles creating outbound SMTP connections
// and sending of SMTP payloads.

package main

import (
	"bufio"
	"fmt"
	"net"
)

func send(env envelope) {
	// TODO: Determine where to send based on RcptTo domain
	conn, err := createNewConnection()

	if err == nil {
		conn.Write([]byte("HELO go-smtpclient\r\n"))
		status, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("<--client ", status)

		conn.Write([]byte("MAIL FROM:"))
		conn.Write([]byte(env.mailFrom))
		conn.Write([]byte("\r\n"))
		status, _ = bufio.NewReader(conn).ReadString('\n')
		fmt.Println("<--client ", status)

		conn.Write([]byte("RCPT TO:"))
		conn.Write([]byte(env.rcptTo))
		conn.Write([]byte("\r\n"))
		status, _ = bufio.NewReader(conn).ReadString('\n')
		fmt.Println("<--client ", status)

		conn.Write([]byte("DATA\r\n"))
		status, _ = bufio.NewReader(conn).ReadString('\n')
		fmt.Println("<--client ", status)

		fmt.Println("SENDING DATA...")
		conn.Write([]byte(env.rawBody))
		conn.Write([]byte("\r\n"))
		status, _ = bufio.NewReader(conn).ReadString('\n')
		fmt.Println("<--client ", status)

	} else {
		fmt.Println("Failed to send message.", err)
		// TODO Don't just drop the message on the floor.
	}

	fmt.Println("Completed SMTP client conversation.")
	conn.Close()
}

func createNewConnection() (conn net.Conn, err error) {

	fmt.Println("Creating new outbound SMTP connection to localhost:2525...")
	conn, err = net.Dial("tcp", "localhost:2525")
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return conn, err
}
