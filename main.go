package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

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

type Peers map[string]*Peer

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

	self := &Peer{
		Address:  cli.ListenAddrress,
		Port:     cli.ListenPort,
		Protocol: "tcp",
	}

	go startListener(cli.ListenAddrress, self)
	time.Sleep(2 * time.Second)

	var newPeers []Peer
	for _, peer := range PeerList {
		Hello(peer, self)
		newPeers := append(newPeers, GetPeers(self, peer)...)
		log.Debug().Str("from", peer.String()).Int("amount", len(newPeers)).Msg("received peers")
	}
	for _, p := range newPeers {
		log.Debug().Str("peer", p.String()).Msg("received peer")
		if _, exists := PeerList[p.Hash()]; !exists && p.Hash() != self.Hash() {
			PeerList[p.Hash()] = &p
			log.Debug().Str("newPeer", p.String()).Msg("received new peer")
		}
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
		address := strings.Split(seed, ":")[0]
		port, err := strconv.Atoi(strings.Split(seed, ":")[1])
		if err != nil {
			log.Error().Err(err).Str("seed", seed).Msg("invalid port for seed")
		} else {
			newPeer := &Peer{Address: address, Port: port, Protocol: "tcp"}
			peers[newPeer.Hash()] = newPeer
			log.Debug().Str("seedNode", newPeer.String()).Msg("added new peer")
		}
	}
	return peers
}

func handleDeadPeer(peer *Peer) {
	log.Info().Str("peer", peer.String()).Msg("removing dead peer")
	delete(PeerList, peer.Hash())
}
