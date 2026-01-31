package peer

import "encoding/json"

const (
	ConnectEnvelope = "CONNECT"
)

type Envelope struct {
	Type    string
	Payload json.RawMessage
}
