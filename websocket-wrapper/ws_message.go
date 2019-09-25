package ww

type WSMessageType int64

const (
	// for business message
	WSMessageTypeBusiness = WSMessageType(1)
	// for control message (close connection)
	WSMessageTypeClose = WSMessageType(2)
)

type WSMessage struct {
	// Type determine the message is business message or not
	Type WSMessageType
	// Data save the concrete datas for translation
	Data []byte
}
