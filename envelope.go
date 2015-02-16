package main

import (
//	"bytes"
	"net"
)

type envelope struct {
	remoteIp net.Addr
	preamble string
	mailFrom string
    rcptTo string
	rawBody string
//	rawBody bytes.Buffer // TODO Make rawBody a stream
}


