package main

import (
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func Hello(peer *Peer, self *Peer) (*Peer, error) {
	msg := &Message{
		PeerID: self.ID,
		Type:   HELLO,
	}

	err := peer.Connect()
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
		return nil, err
	}
	if err := peer.SendMsg(data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
		return nil, err
	}

	raw, err := peer.GetMessage()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	msg = new(Message)
	err = proto.Unmarshal(raw, msg)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return nil, err
	}
	peer.ID = msg.PeerID
	log.Debug().Str("peer", peer.String()).Str("verb", "hello").Msg("")

	myself := PbPeer{
		ID:       self.ID,
		Address:  self.Address,
		Port:     int64(self.Port),
		Protocol: self.Protocol,
	}
	data, err = proto.Marshal(&myself)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
		return nil, err
	}
	if err := peer.SendMsg(data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
		return nil, err
	}

	peer.Connection.Close()

	return peer, nil
}

func handleHello(conn net.Conn, self *Peer, peer *Peer) *Peer {
	msg := &Message{
		PeerID: self.ID,
		Type:   HELLO,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
	}
	if err := peer.SendMsg(data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
	}

	raw, err := peer.GetMessage()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil
	}
	rPeer := new(PbPeer)
	err = proto.Unmarshal(raw, rPeer)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return nil
	}
	peer.ID = rPeer.ID
	peer.Port = int(rPeer.Port)
	peer.Protocol = rPeer.Protocol

	return peer
}

func GetPeers(self *Peer, from *Peer) ([]*Peer, error) {
	if err := from.Connect(); err != nil {
		return nil, err
	}
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

	log.Debug().Str("peer", from.String()).Str("verb", "getpeers").Msg("")
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
		return nil, err
	}
	if err := from.SendMsg(data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
		return nil, err
	}
	raw, err := from.GetMessage()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	from.Connection.Close()

	retrievedPeers := new(PeersList)
	err = proto.Unmarshal(raw, retrievedPeers)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return nil, err
	}

	var peers []*Peer
	for _, peer := range retrievedPeers.Peers {
		p := &Peer{
			ID:       peer.ID,
			Address:  peer.Address,
			Port:     int(peer.Port),
			Protocol: peer.Protocol,
		}
		peers = append(peers, p)
	}

	return peers, nil
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
		sendList.Peers = append(sendList.Peers, &PbPeer{
			ID:       p.ID,
			Address:  p.Address,
			Port:     int64(p.Port),
			Protocol: p.Protocol,
		})
		log.Debug().Str("to", peer.String()).Str("peer", p.String()).Str("id", p.ID).Str("verb", "sentpeers").Msg("")
		if len(sendList.Peers) == limit {
			break
		}
	}

	data, err := proto.Marshal(sendList)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
	}
	if err := peer.SendMsg(data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
	}
}
