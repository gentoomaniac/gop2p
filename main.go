package main

import (
	"os"
	"os/signal"
	"time"

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

type Peers map[[32]byte]*Peer

var (
	PeerList Peers
	RUN      bool
)

var cli struct {
	logging.LoggingConfig

	ListenAddrress string   `short:"l" default:"0.0.0.0" help:"adress to listen on for new connections"`
	ListenPort     int      `short:"p" default:"1234" help:"port to listen on for new connections"`
	ConnectionType string   `short:"t" default:"tcp" enum:"tcp,udp" help:"the type of connection to use"`
	SeedHosts      []string `short:"s" help:"seed host in the form of <address>:<port>"`

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

	PeerList = initialiseSeeds(PeerList, cli.SeedHosts)
	RUN = true
	trapTerm()

	go startListener(cli.ListenAddrress, cli.ListenPort, cli.ConnectionType)
	time.Sleep(2 * time.Second)

	for _, peer := range PeerList {
		log.Debug().Str("seed", peer.String()).Msg("sending hello to seed")
		initHello(peer)
	}

	for RUN {
		time.Sleep(2 * time.Second)
	}
	ctx.Exit(0)
}

func trapTerm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				RUN = false
			}
		}
	}()
}

func initialiseSeeds(peers Peers, seeds []string) Peers {
	peers = make(Peers)

	for _, seed := range seeds {
		newPeer := &Peer{Address: seed, Protocol: "tcp"}
		peers[newPeer.Hash()] = newPeer

		log.Debug().Str("seedNode", newPeer.String()).Msg("added new peer")
	}
	return peers
}

func initHello(peer *Peer) {
	msg := &Message{
		Type:    HELLO,
		Payload: "Hello, seed!",
	}

	err := peer.Connect()
	if err != nil {
		handleDeadPeer(peer)
		return
	}

	log.Debug().Str("peer", peer.String()).Msg("sending hello")
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed Marshalling msg")
	}
	if err := peer.SendMsg(data); err != nil {
		log.Error().Err(err).Msg("failed sending message")
	}

	log.Debug().Str("peer", peer.String()).Msg("waiting hello reply")
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

	log.Debug().Str("payload", msg.Payload).Int("type", int(msg.Type)).Msg("received msg back")

	peer.Connection.Close()
}

func handleDeadPeer(peer *Peer) {
	log.Info().Str("peer", peer.String()).Msg("removing dead peer")
	delete(PeerList, peer.Hash())
}
