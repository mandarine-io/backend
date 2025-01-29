package websocket

type ClientMessage struct {
	ClientID string `json:"clientID"`
	Payload  []byte `json:"payload"`
}

func NewClientMessage(clientID string, payload []byte) ClientMessage {
	return ClientMessage{
		ClientID: clientID,
		Payload:  payload,
	}
}

type BroadcastMessage struct {
	Payload []byte `json:"payload"`
}

func NewBroadcastMessage(payload []byte) BroadcastMessage {
	return BroadcastMessage{
		Payload: payload,
	}
}
