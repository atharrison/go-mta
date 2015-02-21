// Handles creating outbound SMTP connections
// and sending of SMTP payloads.

package main

import (
	"bufio"
	"net"
)

type SmtpClient struct {
	conn net.Conn
	env *envelope
}

func handleSmtpClientConnections(smtpClientChan chan *SmtpClient) {
	Info.Println("SmtpClient Connection Handler Started.")
	for {
		smtpClient := <-smtpClientChan
		Info.Println("Received new SmtpClient, processing outbound SMTP Conversation.")
		sendSmtp(smtpClient.conn, smtpClient.env)
	}
}

func sendSmtp(conn net.Conn, env *envelope) {
	conn.Write([]byte("HELO go-smtpclient\r\n"))
	status, _ := bufio.NewReader(conn).ReadString('\n')
	Debug.Println("<--client ", status)

	conn.Write([]byte("MAIL FROM:"))
	conn.Write([]byte(env.mailFrom))
	conn.Write([]byte("\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	Debug.Println("<--client ", status)

	conn.Write([]byte("RCPT TO:"))
	conn.Write([]byte(env.rcptTo))
	conn.Write([]byte("\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	Debug.Println("<--client ", status)

	conn.Write([]byte("DATA\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	Debug.Println("<--client ", status)

	Debug.Println("SENDING DATA...")
	conn.Write([]byte(env.rawBody))
	conn.Write([]byte("\r\n"))
	status, _ = bufio.NewReader(conn).ReadString('\n')
	Debug.Println("<--client ", status)

	Info.Println("Completed SMTP client conversation.")
	conn.Close()
}
