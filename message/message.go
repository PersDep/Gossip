package message

import (
	"encoding/json"
)

type Message struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Sender int    `json:"sender"`
	Origin int    `json:"origin"`
	Data   string `json:"data"`
}

func (msg Message) ConvertToJsonMsg() []byte {
	buf, error := json.Marshal(msg)

	if error != nil {
		panic(error)
	}

	return buf
}

func ConvertFromJsonMsg(buf []byte) Message {
	msg := &Message{}
	error := json.Unmarshal(buf, msg)

	if error != nil {
		panic(error)
	}

	return *msg
}
