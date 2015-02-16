// Handles accepting new connections
// and managing the incoming conversations.

package main

import (
	"bufio"
	"bytes"
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

	remoteIp := conn.RemoteAddr()
	fmt.Println("Received Connection from [", remoteIp, "]")

	// HELO
	fmt.Println("--> 250 HELO")
	conn.Write([]byte("250 HELO localhost\r\n"))
	status, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("<--", status)

	preamble := strings.Trim(strings.SplitN(status, "HELO ", 2)[1], "\r\n")
	fmt.Println("Preamble [", preamble, "]")

	// MAIL FROM
	conn.Write([]byte("250 OK\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	fmt.Println("<--", status)

	mailFrom := strings.Trim(strings.SplitN(status, "MAIL FROM:", 2)[1], "\r\n")
	fmt.Println("MailFrom:", mailFrom)

    Connection:
	for {

		// RCPT TO
		conn.Write([]byte("250 OK\r\n"))
		status, _ = bufio.NewReader(conn).ReadString('\n')
		fmt.Println("<--", status)

		rcptTo := strings.Trim(strings.SplitN(status, "RCPT TO:", 2)[1], "\r\n")
		fmt.Println("RcptTo:", rcptTo)

		env := envelope{remoteIp, preamble, mailFrom, rcptTo, ""}
		fmt.Println(env)

		// DATA
		conn.Write([]byte("250 OK\r\n"))
		status, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("<--", status)

		conn.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))

		fmt.Println("DATA BLOCK START")
		endOfData := false
		var rawBody bytes.Buffer
		for {
			msg := make([]byte, 256)

			bytesRead, err := conn.Read(msg)
			data := string(msg[:])
			fmt.Println(data)
			rawBody.WriteString(data)

			// Determine if we are at the end, looking for <CR><LF>.<CR><LF>
			lines := strings.Split(data, "\r\n")
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
				break Connection
			}
			if err != nil {
				fmt.Println("Err: ", err)
				break
			}
		}
		fmt.Println("Acknowledging end of DATA.")

		conn.Write([]byte("250 OK\r\n"))
		status, _ = bufio.NewReader(conn).ReadString('\n')
		fmt.Println("<--", status)

		env.rawBody = rawBody.String()

		go handleEnvelope(env)

		if strings.Index(status, "MAIL FROM") == 0 {
			fmt.Println("Detected new message on single connection.")
			mailFrom := strings.Trim(strings.SplitN(status, "MAIL FROM:", 2)[1], "\r\n")
			fmt.Println("MailFrom:", mailFrom)
			continue Connection
		} else if status == "QUIT\r\n" {
			// QUIT
			conn.Write([]byte("221 localhost Service closing transmission channel.\r\n"))
			conn.Close()
			break
		}
	}
}


