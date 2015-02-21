// Handles accepting new SMTP connections
// and managing the incoming conversations.

package main

import (
	"bufio"
	"bytes"
	"net"
	"strings"
)

type SmtpServer struct {
	conn net.Conn
}

func startSmtpServerListener(smtpServerChan chan *SmtpServer) {
	// listen on a port
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		Error.Println(err)
		return
	}
	for {
		// accept a connection
		conn, err := ln.Accept()
		if err != nil {
			Error.Println(err)
			continue
		}
		// handle the connection
		Info.Println("Accepting new Connection, placing on SmtpServer on Channel.")
		server := &SmtpServer{conn}
		smtpServerChan <- server
	}
}

func handleSmtpServerConnections(smtpServerChan chan *SmtpServer, envelopeChan chan *envelope) {
	Info.Println("SmtpServer Connection Handler Started.")
	for {
		// handle an SMTP Conversation
		server := <- smtpServerChan
		Info.Println("Received new SmtpServer, processing inbound SMTP Conversation.")
		receiveSmtp(server.conn, envelopeChan)
	}
}

func receiveSmtp(conn net.Conn, envelopeChan chan *envelope) {

	remoteIp := conn.RemoteAddr()
	Info.Println("Received Connection from [", remoteIp, "]")

	// HELO
	Debug.Println("--> 250 HELO")
	conn.Write([]byte("250 HELO localhost\r\n"))
	status, _ := bufio.NewReader(conn).ReadString('\n')
	Debug.Println("<--", status)

	preamble := strings.Trim(strings.SplitN(status, "HELO ", 2)[1], "\r\n")
	Debug.Println("Preamble [", preamble, "]")

	// MAIL FROM
	conn.Write([]byte("250 OK\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	Debug.Println("<--", status)

	mailFrom := strings.Trim(strings.SplitN(status, "MAIL FROM:", 2)[1], "\r\n")
	Debug.Println("MailFrom:", mailFrom)

Connection:
	for {

		// RCPT TO
		conn.Write([]byte("250 OK\r\n"))
		status, _ = bufio.NewReader(conn).ReadString('\n')
		Debug.Println("<--", status)

		rcptTo := strings.Trim(strings.SplitN(status, "RCPT TO:", 2)[1], "\r\n")
		Debug.Println("RcptTo:", rcptTo)

		env := envelope{remoteIp, preamble, mailFrom, rcptTo, ""}
		Debug.Println(env)

		// DATA
		conn.Write([]byte("250 OK\r\n"))
		status, _ := bufio.NewReader(conn).ReadString('\n')
		Debug.Println("<--", status)

		conn.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))

		Debug.Println("DATA BLOCK START")
		endOfData := false
		var rawBody bytes.Buffer
		for {
			msg := make([]byte, 256)

			bytesRead, err := conn.Read(msg)
			data := string(msg[:])
			Debug.Println(data)
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
				Debug.Println("Received terminating line, done reading DATA")
				break
			}
			if bytesRead == 0 {
				Debug.Println("No bytes read, finished reading DATA.")
				break Connection
			}
			if err != nil {
				Error.Println("Err: ", err)
				break
			}
		}
		Debug.Println("Acknowledging end of DATA.")

		conn.Write([]byte("250 OK\r\n"))
		status, _ = bufio.NewReader(conn).ReadString('\n')
		Debug.Println("<--", status)

		env.rawBody = rawBody.String()

		envelopeChan <- &env

		if strings.Index(status, "MAIL FROM") == 0 {
			Debug.Println("Detected new message on single connection.")
			mailFrom := strings.Trim(strings.SplitN(status, "MAIL FROM:", 2)[1], "\r\n")
			Debug.Println("MailFrom:", mailFrom)
			continue Connection
		} else if status == "QUIT\r\n" {
			// QUIT
			conn.Write([]byte("221 localhost Service closing transmission channel.\r\n"))
			conn.Close()
			break
		}
	}
}
