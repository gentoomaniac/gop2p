package main

import (
	"crypto/sha256"
	"encoding/base64"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/denisbrodbeck/machineid"
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

	RUN = true
	trapTerm()

	self := &Peer{
		Address:  cli.ListenAddrress,
		Port:     cli.ListenPort,
		Protocol: "tcp",
	}
	err := updateHostID(self)
	if err != nil {
		log.Error().Err(err).Msg("could not generate host ID")
		ctx.Exit(1)
	}
	log.Info().Str("hostID", self.ID).Msg("")

	go startListener(cli.ListenAddrress, self)
	time.Sleep(2 * time.Second)

	PeerList = make(Peers)
	seedPeers := initialiseSeeds(cli.SeedHosts)
	var newPeers []*Peer
	for _, peer := range seedPeers {
		peer, err = Hello(peer, self)
		if err != nil {
			handleDeadPeer(peer)
		}

		peers, err := GetPeers(self, peer)
		if err != nil {
			log.Error().Err(err).Str("peer", peer.String()).Msg("error fetching peers")
		} else {
			newPeers = append(newPeers, peers...)
			log.Debug().Str("from", peer.String()).Int("amount", len(newPeers)).Msg("received peers")
		}
	}
	for _, p := range append(newPeers, seedPeers...) {
		if p.ID != self.ID {
			PeerList[p.ID] = p
			log.Debug().Str("peer", p.String()).Str("id", p.ID).Msg("adding peer")
		} else {
			log.Debug().Str("peer", p.String()).Str("id", p.ID).Msg("skipping self")
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

func initialiseSeeds(seeds []string) []*Peer {
	var peers []*Peer

	for _, seed := range seeds {
		address := strings.Split(seed, ":")[0]
		port, err := strconv.Atoi(strings.Split(seed, ":")[1])
		if err != nil {
			log.Error().Err(err).Str("seed", seed).Msg("invalid port for seed")
		} else {
			newPeer := &Peer{Address: address, Port: port, Protocol: "tcp"}
			peers = append(peers, newPeer)
		}
	}
	return peers
}

func handleDeadPeer(peer *Peer) {
	log.Info().Str("peer", peer.String()).Msg("removing dead peer")
	delete(PeerList, peer.ID)
}

func updateHostID(self *Peer) error {
	id, err := machineid.ID()
	if err != nil {
		return err
	}
	h := sha256.Sum256([]byte(self.String() + id))
	self.ID = base64.StdEncoding.EncodeToString(h[:32])

	return nil
}
