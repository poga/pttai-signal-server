package signalserver

import "golang.org/x/crypto/ed25519"

type NodeID ed25519.PublicKey

type Signal struct {
	FromID NodeID
	ToID   NodeID

	Msg   []byte
	Extra []byte
}
