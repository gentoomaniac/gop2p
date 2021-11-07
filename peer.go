package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/rs/zerolog/log"
)

type Peer struct {
	Address    string
	Protocol   string
	Connection net.Conn
}

func (p Peer) String() string {
	return fmt.Sprintf("%s/%s", p.Address, p.Protocol)
}

func (p Peer) Hash() [32]byte {
	return sha256.Sum256([]byte(p.String()))
}

func (p *Peer) Connect() (err error) {
	log.Debug().Str("peer", p.String()).Msg("opened connection")
	p.Connection, err = net.Dial(p.Protocol, p.Address)
	return
}

func (p *Peer) Disconnect() {
	if p.Connection != nil {
		p.Connection.Close()
		log.Debug().Str("peer", p.String()).Msg("closed connection")
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
