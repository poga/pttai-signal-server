package signalserver

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
)

func TestServerIdentifyNodeID(t *testing.T) {
	server := NewServer()

	challenge := server.generateChallenge()

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Errorf("failed: %v", err)
	}

	hash := crypto.Keccak256Hash(challenge)

	// sign challenge with key
	sig := ed25519.Sign(privateKey, hash[:])

	nodeID := NodeID(publicKey)
	resp := &challengeResponse{
		NodeID:    nodeID,
		Hash:      hash,
		Signature: sig,
	}

	err = server.verifyNode(challenge, resp)
	if err != nil {
		t.Errorf("failed: %v", err)
	}
}

func TestServerRemoveFromNodeChannels(t *testing.T) {
	server := NewServer()

	nodeID := NodeID{}

	nodeConn, err := server.newNodeConn(nodeID, nil)
	assert.NoError(t, err)

	_, exists := server.nodeChannels.Load(string(nodeID))
	assert.Equal(t, true, exists)

	server.removeFromNodeChannels(nodeConn)

	_, exists = server.nodeChannels.Load(string(nodeID))
	assert.Equal(t, false, exists)
}
