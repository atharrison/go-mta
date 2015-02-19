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

func main() {
//	Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	runtime.GOMAXPROCS(runtime.NumCPU())
	Info.Println("Starting the Go MTA Server.\n")

	// Start the Server listener
	go server()

	// Die after input is read.
	var input string
	fmt.Scanln(&input)
}
