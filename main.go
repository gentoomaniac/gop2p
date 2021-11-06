package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	"github.com/gentoomaniac/gocli"
	"github.com/gentoomaniac/logging"
)

var (
	version = "unknown"
	commit  = "unknown"
	binName = "unknown"
	builtBy = "unknown"
	date    = "unknown"
)

var cli struct {
	logging.LoggingConfig

	ListenAddrress string   `short:"l" default:"0.0.0.0" help:"adress to listen on for new connections"`
	ListenPort     int      `short:"p" default:"1234" help:"port to listen on for new connections"`
	ConnectionType string   `short:"t" default:"tcp" enum:"tcp,udp" help:"the type of connection to use"`
	SeedHosts      []string `short:"s" help:"seed host and port"`

	Version gocli.VersionFlag `short:"V" help:"Display version."`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Vars{
		"version": version,
		"commit":  commit,
		"binName": binName,
		"builtBy": builtBy,
		"date":    date,
	})
	logging.Setup(&cli.LoggingConfig)

	go startListener(cli.ListenAddrress, cli.ListenPort, cli.ConnectionType)

	for _, seed := range cli.SeedHosts {
		s := strings.Split(seed, ":")
		host := s[0]
		port, err := strconv.Atoi(s[1])
		if err != nil {
			log.Error().Err(err).Msg("")
			ctx.Exit(1)
		}

		msg := &Message{
			Type:    HELLO,
			Payload: "Hello, Protobuf!",
		}
		msg2 := &NewPeer{
			Address:  "10.1.1.121",
			Port:     1236,
			Protocol: "tcp",
		}

		peer := &Peer{Address: host, Port: port, ConnectionType: "tcp"}
		log.Debug().Str("seedNode", peer.String()).Msg("sending message")
		peer.Connect()

		var data []byte
		if data, err = proto.Marshal(msg); err != nil {
			log.Error().Err(err).Msg("failed Marshalling msg")
		}
		if err := peer.SendMsg(data); err != nil {
			log.Error().Err(err).Msg("failed sending message")
		}

		if data, err = proto.Marshal(msg2); err != nil {
			log.Error().Err(err).Msg("failed Marshalling msg")
		}
		if err := peer.SendMsg(data); err != nil {
			log.Error().Err(err).Msg("failed sending message")
		}
		peer.connection.Close()
	}

	fmt.Scanln()
	ctx.Exit(0)
}
