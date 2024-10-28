package websocket

type ClientMessage struct {
	ClientId string `json:"clientId"`
	Payload  []byte `json:"payload"`
}

func NewClientMessage(clientId string, payload []byte) ClientMessage {
	return ClientMessage{
		ClientId: clientId,
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
