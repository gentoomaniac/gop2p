package main

import (
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func startListener(host string, port int, connType string) {
	l, err := net.Listen(connType, host+":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Error().Err(err).Msg("failed to start listener")
		os.Exit(1)
	}
	defer l.Close()

	log.Info().Str("listenAddress", host).Int("listenPort", port).Str("connectionType", connType).Msg("starting to listen for new connections ...")
	var handlers sync.WaitGroup
	for RUN {
		conn, err := l.Accept()
		if err != nil {
			log.Error().Err(err).Msg("could not accept message")
			os.Exit(1)
		}
		handlers.Add(1)
		go handleIncomingMessage(conn, &handlers)
	}
}

func handleIncomingMessage(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	s := strings.Split(conn.RemoteAddr().String(), ":")
	address := strings.Join(s[:len(s)-1], ":")
	port, _ := strconv.Atoi(s[len(s)-1])
	peer := &Peer{
		Address:    address,
		Port:       port,
		Protocol:   "tcp",
		Connection: conn,
	}
	log.Debug().Str("peer", peer.String()).Msg("new peer connected")

	raw, err := peer.GetMessage()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	msg := new(Message)
	err = proto.Unmarshal(raw, msg)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return
	}
	log.Debug().Str("payload", msg.Payload).Int("type", int(msg.Type)).Msg("received msg back")

	switch msg.GetType() {
	case HELLO:
		handleHello(conn, peer, PeerList)

	default:
		log.Warn().Str("message", msg.GetPayload()).Int64("type", msg.GetType()).Msg("got unhandled message type")
	}

	log.Debug().Str("ip", conn.RemoteAddr().String()).Msg("peer finished")
}

func handleHello(conn net.Conn, peer *Peer, peers Peers) {
	msg := &Message{
		Type:    HELLO,
		Payload: "hello, peer!",
	}

	log.Debug().Str("peer", peer.String()).Msg("sending hello back")
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
