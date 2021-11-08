package main

import (
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func Hello(peer *Peer, self *Peer) {
	msg := &Message{
		Type: HELLO,
	}

	err := peer.Connect()
	if err != nil {
		handleDeadPeer(peer)
		return
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
		return
	}

	msg = new(Message)
	err = proto.Unmarshal(raw, msg)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return
	}
	log.Debug().Str("peer", peer.String()).Str("verb", "hello").Msg("")

	myself := PbPeer{
		Address:  self.Address,
		Port:     int64(self.Port),
		Protocol: self.Protocol,
	}
	data, err = proto.Marshal(&myself)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
	}
	if err := peer.SendMsg(data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
	}

	peer.Connection.Close()
}
func handleHello(conn net.Conn, peer *Peer, peers Peers) {
	msg := &Message{
		Type: HELLO,
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
		return
	}
	rPeer := new(PbPeer)
	err = proto.Unmarshal(raw, rPeer)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return
	}
	peer.Port = int(rPeer.Port)
	peer.Protocol = rPeer.Protocol

	if _, exists := peers[peer.Hash()]; !exists {
		peers[peer.Hash()] = peer
		log.Info().Str("newPeer", peer.String()).Msg("got new peer")
	}

}

func GetPeers(self *Peer, from *Peer) []Peer {
	from.Connect()
	options := &GetPeersOptions{
		Limit: 0,
	}
	optionsBytes, err := proto.Marshal(options)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
		return nil
	}

	msg := &Message{
		Type:    GETPEERS,
		Payload: optionsBytes,
	}

	log.Debug().Str("peer", from.String()).Str("verb", "getpeers").Msg("")
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
	}
	if err := from.SendMsg(data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
	}
	raw, err := from.GetMessage()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil
	}
	from.Connection.Close()

	retrievedPeers := new(PeersList)
	err = proto.Unmarshal(raw, retrievedPeers)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return nil
	}

	var peers []Peer
	for _, peer := range retrievedPeers.Peers {
		p := Peer{
			Address:  peer.Address,
			Port:     int(peer.Port),
			Protocol: peer.Protocol,
		}
		if !self.Equal(p) {
			peers = append(peers, p)
		}
	}

	return peers
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
			Address:  p.Address,
			Port:     int64(p.Port),
			Protocol: p.Protocol,
		})
		log.Debug().Str("to", peer.String()).Str("peer", p.String()).Str("verb", "sentpeers").Msg("")
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
