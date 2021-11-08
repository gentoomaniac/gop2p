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

func startListener(host string, self *Peer) {
	l, err := net.Listen(self.Protocol, self.Address+":"+strconv.Itoa(int(self.Port)))
	if err != nil {
		log.Error().Err(err).Msg("failed to start listener")
		os.Exit(1)
	}
	defer l.Close()

	log.Info().Str("listenAddress", host).Int("listenPort", self.Port).Str("protocol", self.Protocol).Msg("starting to listen for new connections ...")
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

	switch msg.GetType() {
	case HELLO:
		log.Debug().Str("peer", peer.String()).Str("verb", "hello").Msg("")
		handleHello(conn, peer, PeerList)
	case GETPEERS:
		log.Debug().Str("peer", peer.String()).Str("verb", "getpeers").Msg("")
		handleGetPeers(conn, peer, PeerList, msg)

	default:
		log.Warn().Str("message", string(msg.Payload)).Int64("type", msg.GetType()).Msg("got unhandled message type")
	}
}
