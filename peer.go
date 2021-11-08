package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

type Peer struct {
	ID         string
	Address    string
	Port       int
	Protocol   string
	Connection net.Conn
}

func (p Peer) String() string {
	return fmt.Sprintf("%s:%d/%s", p.Address, p.Port, p.Protocol)
}

func (p Peer) Equal(other Peer) bool {
	if p.Address != other.Address || p.Port != other.Port || p.Protocol != other.Protocol {
		return false
	}

	return true
}

func (p *Peer) Connect() (err error) {
	p.Connection, err = net.Dial(p.Protocol, p.Address+":"+strconv.Itoa(p.Port))
	return
}

func (p *Peer) Disconnect() {
	if p.Connection != nil {
		p.Connection.Close()
	}
}

func (p *Peer) SendMsg(msg []byte) error {
	var b [16]byte
	bs := b[:16]
	binary.LittleEndian.PutUint64(bs, uint64(len(msg)))
	_, err := p.Connection.Write(bs)
	if err != nil {
		return err
	}
	_, err = p.Connection.Write(msg)

	return err
}

func (p *Peer) GetMessage() ([]byte, error) {
	var b [16]byte
	bs := b[:16]

	_, err := io.ReadFull(p.Connection, bs)
	if err != nil {
		return nil, err
	}
	numBytes := uint64(binary.LittleEndian.Uint64(bs))
	data := make([]byte, numBytes)
	_, err = io.ReadFull(p.Connection, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
