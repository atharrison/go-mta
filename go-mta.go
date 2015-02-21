package main

import (
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"os"
	"runtime"
)

func Init(

	// Build Loggers
    traceHandle io.Writer,
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
				log.Ldate|log.Ltime|log.Lshortfile)

	Debug = log.New(debugHandle,
		"TRACE: ",
				log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
				log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
				log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
				log.Ldate|log.Ltime|log.Lshortfile)

}

const SmtpServerConnectionCount = 30
const SmtpClientConnectionCount = 30
const DispatcherThreads = 30

func main() {
//	Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	// Create channels
	smtpServerChan := make(chan *SmtpServer)
	envelopeChan := make(chan *envelope)
	smtpClientChan := make(chan *SmtpClient)

	runtime.GOMAXPROCS(runtime.NumCPU())
	Info.Println("Starting the Go MTA Server.\n")

	// Start the Server listener
	go startSmtpServerListener(smtpServerChan)

	// Handle new server connections
	for i := 0; i < SmtpServerConnectionCount; i++ {
		go handleSmtpServerConnections(smtpServerChan, envelopeChan)
	}
	for i := 0; i < DispatcherThreads; i++ {
		go handleDispatcher(envelopeChan, smtpClientChan)
	}
	for i := 0; i < SmtpClientConnectionCount; i++ {
		go handleSmtpClientConnections(smtpClientChan)
	}

	// Die after input is read.
	var input string
	fmt.Scanln(&input)
}
