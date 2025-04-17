package common

type MsgRsp struct {
	Type        string `json:"message_type"`
	UnReadCount int    `json:"unread_count"`
	Payload     any    `json:"payload"`
}
