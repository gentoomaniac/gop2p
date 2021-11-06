package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func startListener(host string, port int, connType string) {
	l, err := net.Listen(connType, host+":"+strconv.Itoa(port))
	if err != nil {
		log.Error().Err(err).Msg("failed to start listener")
		os.Exit(1)
	}
	defer l.Close()

	log.Info().Str("listenAddress", host).Int("listenPort", port).Str("connectionType", connType).Msg("starting to listen for new connections ...")
	var handlers sync.WaitGroup
	for {
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

	log.Debug().Str("ip", conn.RemoteAddr().String()).Msg("new peer connected")

	var buf bytes.Buffer
	_, err := io.Copy(&buf, conn)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	msg := new(Message)
	err = proto.Unmarshal(buf.Bytes(), msg)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return
	}
	log.Info().Str("message", msg.GetPayload()).Int64("type", msg.GetType()).Msg("")

	peer := new(NewPeer)
	err = proto.Unmarshal(buf.Bytes(), peer)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshalling message")
		return
	}
	log.Info().Str("newPeer", fmt.Sprintf("%s:%d/%s", peer.GetAddress(), peer.GetPort(), peer.GetProtocol())).Msg("got new peer")

	log.Debug().Str("ip", conn.RemoteAddr().String()).Msg("peer finished")
}
