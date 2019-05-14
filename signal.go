package signalserver

import "golang.org/x/crypto/ed25519"

type NodeID ed25519.PublicKey

func (id NodeID) Equal(x NodeID) bool {
	return string(id) == string(x)
}

type Signal struct {
	FromID NodeID
	ToID   NodeID

	Msg   []byte
	Extra []byte
}
