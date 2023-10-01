package plexer

import (
	"bufio"
	"net"
)

type Plexer struct {
	Request  *bufio.Reader
	Response *bufio.Writer

	ReqConn net.Conn
	ResConn net.Conn
}

func (p *Plexer) Close() {
	p.ReqConn.Close()
	p.ResConn.Close()
}
