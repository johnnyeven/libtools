package ww

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WSServerInterface interface {
	OnMessage(WSHandle, []byte, error)
	OnClose(WSHandle)
}

type WSServerHandler struct {
	On WSServerInterface
}

// OnHandler contains a connection and open a routine for recv
func (ws *WSServerHandler) OnHandler(w http.ResponseWriter, r *http.Request) {
	// check ,if the interface is nil
	if ws.On == nil {
		panic("WSServerHandler ")
	}
	// use the default options
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Error("OnHandler,upgrade the general request connection to a websocket connection error:%s", err.Error())
		return
	}
	defer c.Close()
	handle := wss.saveConn(c)
	writeChan := make(chan *WSMessage, 10)
	wss.saveChan(handle, writeChan)
	go func() {
		closeMsg := WSMessage{Type: WSMessageTypeClose}
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				remoteAddr := c.RemoteAddr().String()
				Error("OnHandler,ReadMessage error:%s,handle:%d,remote addr:%s,message type:%d",
					err.Error(), handle, remoteAddr, mt)
				ws.On.OnMessage(handle, nil, err)
				// additional send close msg
				writeChan <- &closeMsg
				break
			}
			switch mt {
			case websocket.TextMessage:
				ws.On.OnMessage(handle, message, nil)
			case websocket.BinaryMessage:
				//Debug("recv keep-alive message")
				// do nothing,binray message only for keep-alive message
			default:
				remoteAddr := c.RemoteAddr().String()
				Error("OnHandler,ReadMessage handle:%d,remote addr:%s,unsupport message type:%d",
					handle, remoteAddr, mt)
				// additional send close msg
				writeChan <- &closeMsg
			}
		}
	}()
WriteDone:
	for {
		select {
		case msg, ok := <-writeChan:
			if !ok {
				break WriteDone
			}
			switch {
			case msg.Type == WSMessageTypeBusiness:
				err = c.WriteMessage(websocket.TextMessage, msg.Data)
				if err != nil {
					remoteAddr := c.RemoteAddr().String()
					Error("OnHandler,WriteMessage error:%s,handle:%d,remote addr:%s,msg type:%d",
						err.Error(), handle, remoteAddr, msg.Type)
					break WriteDone
				}
			case msg.Type == WSMessageTypeClose:
				err := c.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					remoteAddr := c.RemoteAddr().String()
					Error("OnHandler,WriteMessage error:%s,handle:%d,remote addr:%s,msg type:%d",
						err.Error(), handle, remoteAddr, msg.Type)
				}
				break WriteDone
			}
		}
	}
	// clean resources
	close(writeChan)
	wss.removeChan(handle)
	wss.removeConn(c)
	ws.On.OnClose(handle)
}
