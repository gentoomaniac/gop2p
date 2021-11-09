package main

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

const (
	HELLO int64 = iota + 1
	GETPEERS
)

func SendMsg(conn net.Conn, msg []byte) error {
	var b [16]byte
	bs := b[:16]
	binary.LittleEndian.PutUint64(bs, uint64(len(msg)))
	_, err := conn.Write(bs)
	if err != nil {
		return err
	}
	_, err = conn.Write(msg)

	return err
}

func GetMessage(conn net.Conn) ([]byte, error) {
	var b [16]byte
	bs := b[:16]

	_, err := io.ReadFull(conn, bs)
	if err != nil {
		return nil, err
	}
	numBytes := uint64(binary.LittleEndian.Uint64(bs))
	data := make([]byte, numBytes)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Hello(peer *Peer, self *Peer) (*Peer, error) {
	msg := &Message{
		PeerID: self.Id,
		Type:   HELLO,
	}

	conn, err := net.Dial(peer.Protocol, peer.Address+":"+strconv.Itoa(int(peer.Port)))
	if err != nil {
		return peer, err
	}
	defer conn.Close()

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
		return peer, err
	}
	if err := SendMsg(conn, data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
		return peer, err
	}

	raw, err := GetMessage(conn)
	if err != nil {
		log.Error().Err(err).Msg("")
		return peer, err
	}

	msg = new(Message)
	err = proto.Unmarshal(raw, msg)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return peer, err
	}
	peer.Id = msg.PeerID
	log.Debug().Str("peer", peer.ConnectString()).Str("verb", "hello").Msg("")

	data, err = proto.Marshal(self)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
		return peer, err
	}
	if err := SendMsg(conn, data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
		return peer, err
	}

	conn.Close()
	return peer, nil
}

func handleHello(conn net.Conn, self *Peer, peer *Peer) *Peer {
	msg := &Message{
		PeerID: self.Id,
		Type:   HELLO,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
	}
	if err := SendMsg(conn, data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
	}

	raw, err := GetMessage(conn)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil
	}
	rPeer := new(Peer)
	err = proto.Unmarshal(raw, rPeer)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return nil
	}
	rPeer.Address = peer.Address

	return rPeer
}

func GetPeers(self *Peer, from *Peer) ([]*Peer, error) {
	conn, err := net.Dial(from.Protocol, from.Address+":"+strconv.Itoa(int(from.Port)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	options := &GetPeersOptions{
		Limit: 0,
	}
	optionsBytes, err := proto.Marshal(options)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
		return nil, err
	}

	msg := &Message{
		Type:    GETPEERS,
		Payload: optionsBytes,
	}

	log.Debug().Str("peer", from.ConnectString()).Str("verb", "getpeers").Msg("")
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
		return nil, err
	}
	if err := SendMsg(conn, data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
		return nil, err
	}
	raw, err := GetMessage(conn)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	conn.Close()

	retrievedPeers := new(PeersList)
	err = proto.Unmarshal(raw, retrievedPeers)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return nil, err
	}

	return retrievedPeers.Peers, nil
}

func handleGetPeers(conn net.Conn, peer *Peer, peers Peers, msg *Message) {
	opts := new(GetPeersOptions)

	err := proto.Unmarshal(msg.Payload, opts)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return
	}

	limit := len(peers)
	if int(opts.Limit) < len(peers) {
		limit = int(opts.Limit)
	}

	sendList := new(PeersList)
	for _, p := range peers {
		sendList.Peers = append(sendList.Peers, p)
		log.Debug().Str("to", peer.ConnectString()).Str("peer", p.ConnectString()).Str("id", p.Id).Str("verb", "sentpeers").Msg("")
		if len(sendList.Peers) == limit {
			break
		}
	}

	data, err := proto.Marshal(sendList)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
	}
	if err := SendMsg(conn, data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
	}
}
