package peer

import "encoding/json"

const (
	ConnectEnvelope = "CONNECT"
	MsgEnvelope     = "MSG"
)

type Envelope struct {
	Type    string
	Payload json.RawMessage
}
