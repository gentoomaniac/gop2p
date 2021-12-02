package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
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

type Peers map[[32]byte]*Peer

var (
	PeerList Peers
	RUN      bool
)

var cli struct {
	logging.LoggingConfig

	ListenAddrress string   `short:"l" default:"0.0.0.0" help:"adress to listen on for new connections"`
	ListenPort     int64    `short:"p" default:"1234" help:"port to listen on for new connections"`
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
	id, err := genHostID(self)
	if err != nil {
		log.Error().Err(err).Msg("could not generate host ID")
		ctx.Exit(1)
	}
	self.Id = id[:32]
	log.Info().Str("hostID", self.IDString()).Msg("")

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
			log.Error().Err(err).Str("peer", peer.ConnectString()).Msg("error fetching peers")
		} else {
			newPeers = append(newPeers, peers...)
			log.Debug().Str("from", peer.ConnectString()).Int("amount", len(newPeers)).Msg("received peers")
		}
	}
	for _, p := range newPeers {
		if !bytes.Equal(p.Id, self.Id) {
			if _, err := Hello(p, self); err == nil {
				PeerList[p.Sha256()] = p
				log.Debug().Str("peer", p.ConnectString()).Str("id", p.IDString()).Str("distance", self.Distance(p.Id).String()).Msg("adding peer")
			} else {
				log.Debug().Str("peer", p.ConnectString()).Str("id", p.IDString()).Msg("dropping dead peer")
			}
		} else {
			log.Debug().Str("peer", p.ConnectString()).Str("id", p.IDString()).Msg("skipping self")
		}
	}
	for _, p := range seedPeers {
		PeerList[p.Sha256()] = p
		log.Debug().Str("peer", p.ConnectString()).Str("id", p.IDString()).Str("distance", self.Distance(p.Id).String()).Msg("adding seed peer")
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
			newPeer := &Peer{Address: address, Port: int64(port), Protocol: "tcp"}
			peers = append(peers, newPeer)
		}
	}
	return peers
}

func handleDeadPeer(peer *Peer) {
	log.Info().Str("peer", peer.ConnectString()).Msg("removing dead peer")
	if PeerList != nil {
		delete(PeerList, peer.Sha256())
	}
}

func genHostID(self *Peer) ([32]byte, error) {
	id, err := machineid.ID()
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256([]byte(self.ConnectString() + id)), nil
}

func (p *Peer) ConnectString() string {
	return p.Address + ":" + strconv.Itoa(int(p.Port)) + "/" + p.Protocol
}

func (p *Peer) Sha256() [32]byte {
	sha := [32]byte{}
	copy(sha[:], p.Id[:32])

	return sha
}

func (p *Peer) IDString() string {
	return hex.EncodeToString(p.Id)
}

func (p *Peer) Distance(toID []byte) *big.Int {
	to := &big.Int{}
	to.SetBytes(toID)

	from := &big.Int{}
	from.SetBytes(p.Id)

	distance := &big.Int{}
	distance.Xor(from, to)
	return distance
}
