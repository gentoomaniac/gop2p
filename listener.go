package main

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
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

	log.Debug().Str("ip", conn.RemoteAddr().String()).Msg("new peer connected")

	connbuf := bufio.NewReader(conn)
	for {
		str, err := connbuf.ReadString('\n')
		if err != nil {
			break
		}

		if len(str) > 0 {
			log.Debug().Str("msg", str).Msg("msg received")
		}
	}

	log.Debug().Str("ip", conn.RemoteAddr().String()).Msg("peer finished")
}
