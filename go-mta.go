package main

import (
	"log"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
		"DEBUG: ",
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

	service := NewServiceHandler()

	// Start the Server listener
	cl := NewConnectionListener(smtpServerChan)
	go cl.start()

	// Handle new server connections
	for i := 0; i < SmtpServerConnectionCount; i++ {
		go service.handleSmtpServerConnections(smtpServerChan, envelopeChan)
	}
	for i := 0; i < DispatcherThreads; i++ {
		go handleDispatcher(envelopeChan, smtpClientChan)
	}
	for i := 0; i < SmtpClientConnectionCount; i++ {
		go handleSmtpClientConnections(smtpClientChan)
	}

	// Handle SIGINT and SIGTERM.
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	Info.Println(<-signalCh)

	// Stop the service gracefully.
	service.Stop()
}
