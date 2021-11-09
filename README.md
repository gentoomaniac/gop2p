# p2p Experiments

## tl;dr

This project is for educational purposes and tries to implement a basic p2p network using protobuf for message encoding.

## Status

Peers can exchange information about other known peers.

## Protocol

### Wire

* a little endian encoded uint64 that specifies the size of the next message
* the message itself

### Application level

#### HELLO

A peer sends a `Message` with type `HELLO` to announce itself to another peer.

The other node is replying with a `HELLO`.

The initiating peer is replying with a `Peer` containing it's own data
(ToDo: move this into the payload of the `HELLO`)


#### GETPEERS

A peer sends a `Message` with type `GETPEERS` to ask for other known peers.
The payload of the request is `GetPeersOptions` object.

The receiving node will reply with `[n <=limit]` `Peer` objects.
