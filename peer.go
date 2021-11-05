package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Peer struct {
	Address        string
	Port           int
	ConnectionType string
	connection     net.Conn
}

func (p Peer) String() string {
	return fmt.Sprintf("%s:%d/%s", p.Address, p.Port, p.ConnectionType)
}

func (p *Peer) Connect() (err error) {
	log.Debug().Str("peer", p.String()).Msg("opened connection")
	p.connection, err = net.Dial(p.ConnectionType, p.Address+":"+strconv.Itoa(p.Port))
	return
}

func (p *Peer) Disconnect() {
	if p.connection != nil {
		p.connection.Close()
		log.Debug().Str("peer", p.String()).Msg("closed connection")
	}
}

func (p *Peer) SendMsg(msg []byte) {
	p.connection.Write(append(msg, '\n'))
}
