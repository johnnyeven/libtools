package ww

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// WSHandle means a websocket connection
type WSHandle uint64

// wsServer save websocket server connection info
type wsServer struct {
	// the pair of map save relate info,but the key of connMap is the value of connHandleMap
	connMap       map[*websocket.Conn]WSHandle
	connHandleMap map[WSHandle]*websocket.Conn
	connChanMap   map[WSHandle]chan *WSMessage
	// the rwlock for concurrent safety
	rwLock sync.RWMutex
}

var (
	// connHandle keep the conn has different handle value
	connHandle uint64
	// wss keep server's info
	wss wsServer
)

func init() {
	// init the wss fields
	wss = wsServer{}
	wss.rwLock.Lock()
	defer wss.rwLock.Unlock()
	if wss.connMap == nil {
		wss.connMap = make(map[*websocket.Conn]WSHandle)
	}
	if wss.connHandleMap == nil {
		wss.connHandleMap = make(map[WSHandle]*websocket.Conn)
	}
	if wss.connChanMap == nil {
		wss.connChanMap = make(map[WSHandle]chan *WSMessage)
	}
	// add the other value for initialize
}

// new a WSHandle as the index of conn
func (wss *wsServer) newWSHandle() WSHandle {
	u := atomic.AddUint64(&connHandle, 1)
	return WSHandle(u)
}

// save a conn
func (wss *wsServer) saveConn(c *websocket.Conn) WSHandle {
	hd := wss.newWSHandle()
	wss.rwLock.Lock()
	defer wss.rwLock.Unlock()
	wss.connMap[c] = hd
	wss.connHandleMap[hd] = c
	return hd
}

// remove a conn
func (wss *wsServer) removeConnByIndex(hd WSHandle) {
	wss.rwLock.Lock()
	defer wss.rwLock.Unlock()
	if k, ok := wss.connHandleMap[hd]; ok {
		delete(wss.connMap, k)
		delete(wss.connHandleMap, hd)
	}
}

// remove a conn by conn
func (wss *wsServer) removeConn(conn *websocket.Conn) {
	wss.rwLock.Lock()
	defer wss.rwLock.Unlock()
	if k, ok := wss.connMap[conn]; ok {
		delete(wss.connHandleMap, k)
		delete(wss.connMap, conn)
	}
}

// save chan
func (wss *wsServer) saveChan(hd WSHandle, c chan *WSMessage) {
	wss.rwLock.Lock()
	defer wss.rwLock.Unlock()
	wss.connChanMap[hd] = c
}

// remove chan
func (wss *wsServer) removeChan(hd WSHandle) {
	wss.rwLock.Lock()
	defer wss.rwLock.Unlock()
	delete(wss.connChanMap, hd)
}

// get chan
func (wss *wsServer) getChan(hd WSHandle) (chan *WSMessage, error) {
	wss.rwLock.RLock()
	defer wss.rwLock.RUnlock()
	c, ok := wss.connChanMap[hd]
	if !ok {
		errorMsg := fmt.Sprintf("handle:%d not found chan", hd)
		return nil, errors.New(errorMsg)
	}
	return c, nil
}

// SendMessge ,send a message to chan
func (h WSHandle) SendMessage(msg *WSMessage) (err error) {
	if c, err := wss.getChan(h); err != nil {
		fmt.Println("SendMessge error:", err.Error())
		return err
	} else {
		defer func() {
			if e := recover(); e != nil {
				Error("SendMessge recover error:%+v", e)
				err = errors.New("SendMessge erros,write to chan fail")
			}
		}()
		c <- msg
	}
	return nil
}
