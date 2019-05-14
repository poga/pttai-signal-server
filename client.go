package signalserver

import (
	"net/url"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

type Client struct {
	nodeID NodeID
	Conn   *Conn
}

/*
Send sends the normal-msg to toID.
*/
func (c *Client) Send(toID NodeID, msg []byte, extra []byte) error {
	signal := &Signal{FromID: c.nodeID, ToID: toID, Msg: msg, Extra: extra}

	err := c.Conn.WsConn.WriteJSON(signal)
	if err != nil {
		return err
	}

	return nil
}

/*
Receive receives the normal-msg from signal-server.
*/
func (c *Client) Receive() (*Signal, error) {
	signal := &Signal{}
	err := c.Conn.WsConn.ReadJSON(signal)
	if err != nil {
		return nil, err
	}

	if !sameNodeID(c.nodeID, signal.ToID) {
		return nil, ErrInvalidNodeID
	}

	return signal, nil
}

/*
NewClient init a new client and pass the challenge from the signal-server.
*/
func NewClient(nodeID NodeID, privKey *ed25519.PrivateKey, url url.URL) (*Client, error) {
	wsConn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Dial failed")
	}

	c := &challenge{}
	err = wsConn.ReadJSON(c)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse challenge")
	}

	resp, err := respondChallenge(nodeID, privKey, c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to respond challenge")
	}

	err = wsConn.WriteJSON(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write challenge response")
	}

	cack := &challengeAck{}
	err = wsConn.ReadJSON(cack)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read ack")
	}

	if !sameNodeID(cack.NodeID, nodeID) {
		return nil, ErrInvalidNodeID
	}

	conn := &Conn{isClosed: 0, WsConn: wsConn}

	return &Client{nodeID: nodeID, Conn: conn}, nil
}

func respondChallenge(nodeID NodeID, privKey *ed25519.PrivateKey, c *challenge) (*challengeResponse, error) {
	hash := crypto.Keccak256Hash(c.Challenge)

	sig := ed25519.Sign(*privKey, hash[:])

	challengeResponse := &challengeResponse{NodeID: nodeID, Signature: sig, Hash: hash}

	return challengeResponse, nil
}

func (c *Client) Close() {
	c.Conn.Close()
}

func sameNodeID(a NodeID, b NodeID) bool {
	return string(a) == string(b)
}
