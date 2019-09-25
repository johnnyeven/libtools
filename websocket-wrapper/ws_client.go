package ww

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	// remote server's listen addr, like "localhost:8080"
	remoteAddr string
	// remote server's path, like "/echo"
	remotePath string
	// reconnect time
	reconnectTime time.Duration
	// on message func
	onMessage func([]byte, error)
	// on close func
	onClose func()
	// client send message chan
	c chan *WSMessage
	// conn
	conn *websocket.Conn
	// force quit
	isForceQuit bool
	// lock for concurrent safety
	rwLock sync.RWMutex
	// reconnectMsg,when reconnect or connect send it
	reconnectMsg *WSMessage
}

func (ws *WSClient) connectRemote() error {
	u := url.URL{Scheme: "ws", Host: ws.remoteAddr, Path: ws.remotePath}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		Error("dial err:%s,addr:%s", err.Error(), u.String())
		return err
	}
	ws.rwLock.Lock()
	defer ws.rwLock.Unlock()
	ws.conn = c
	return nil
}

func (ws *WSClient) disconnectRemote() error {
	ws.rwLock.Lock()
	defer ws.rwLock.Unlock()
	//TODO send close message first
	if ws.c != nil {
		close(ws.c)
		ws.c = nil
	}
	if ws.conn != nil {
		ws.conn.Close()
		ws.conn = nil
	}
	return nil
}

func (ws *WSClient) Shutdown() error {
	closeMsg := WSMessage{
		Type: WSMessageTypeClose,
	}
	ws.SendMessge(&closeMsg)
	// need set forceQuit true
	ws.isForceQuit = true
	time.Sleep(time.Second)
	return nil
}

func (ws *WSClient) Startup(remoteAddr, remotePath string, reconnectTime time.Duration,
	onMessage func([]byte, error), onClose func(), reconnectMsg *WSMessage) error {
	// check the input parameters
	if len(remoteAddr) == 0 || len(remotePath) == 0 || nil == onMessage || nil == onClose {
		errStr := fmt.Sprintf("invalid input paras,remoteAddr[%s],remotePath[%s],onMessage[%p],onClose[%p]",
			remoteAddr, remotePath, onMessage, onClose)
		Error(errStr)
		return errors.New(errStr)
	}
	ws.remoteAddr = remoteAddr
	ws.remotePath = remotePath
	ws.reconnectTime = reconnectTime
	ws.onMessage = onMessage
	ws.onClose = onClose
	{
		ws.rwLock.Lock()
		ws.reconnectMsg = reconnectMsg
		ws.rwLock.Unlock()
	}
	// if the first connect is failed,caller need handle the error
	if err := ws.connectRemote(); err != nil {
		Error("WSClient,connect remote server error:%s", err.Error())
		return err
	}
	// the below anonymous go routine for handle client's recv,send,reconnect and quit
	go ws.run()
	// sleep 10 ms ,make sure the client's inner chan has already created
	time.Sleep(10 * time.Microsecond)

	return nil
}

func (ws *WSClient) SendMessge(msg *WSMessage) (err error) {
	if msg == nil {
		return nil
	}
	ws.rwLock.RLock()
	defer ws.rwLock.RUnlock()
	if ws.c == nil {
		Error("WSClient SendMessge chan is nil")
		return errors.New("SendMessge,the chan is nil")
	}
	ws.c <- msg
	defer func() {
		if e := recover(); e != nil {
			Error("SendMessge recover error:%+v", e)
			err = errors.New("SendMessge erros,write to chan fail")
		}
	}()
	return nil
}

func (ws *WSClient) run() {
	for {
		{
			ws.rwLock.Lock()
			ws.c = make(chan *WSMessage, 10)
			ws.rwLock.Unlock()
		}
		// the below anonymous go routine for handle the the client's recv
		reconnectChan := make(chan struct{})
		go func() {
			for {
				mt, message, err := ws.conn.ReadMessage()
				if err != nil {
					Error("WSClient,ReadMessage error:%s", err.Error())
					ws.onMessage(nil, err)
					close(reconnectChan)
					break
				}
				switch mt {
				case websocket.TextMessage:
					ws.onMessage(message, nil)
				case websocket.CloseMessage:
					ws.onClose()
					close(reconnectChan)
					break
				default:
					errMsg := fmt.Sprintf("websocket client unsupport message type:%d", mt)
					panic(errors.New(errMsg))
				}
			}
		}()
		{
			ws.rwLock.RLock()
			// if the ws.reconnectMsg is not nil,send it
			if nil != ws.reconnectMsg {
				// for thread safe,send the copy
				copyMsg := *ws.reconnectMsg
				ws.SendMessge(&copyMsg)
			}
			ws.rwLock.RUnlock()
		}
		// the below loop for handle the client's send
	WriteDone:
		for {
			keepAliveTime := 30 * time.Second
			keepAliveTimer := time.NewTimer(keepAliveTime)
			select {
			case msg, ok := <-ws.c:
				if !ok {
					break WriteDone
				}
				switch {
				case msg.Type == WSMessageTypeBusiness:
					err := ws.conn.WriteMessage(websocket.TextMessage, msg.Data)
					if err != nil {
						remoteAddr := ws.conn.RemoteAddr().String()
						Error("WSClient,WriteMessage error:%s,remote addr:%s,msg type:%d",
							err.Error(), remoteAddr, msg.Type)
						break WriteDone
					}
					keepAliveTimer.Reset(keepAliveTime)
				case msg.Type == WSMessageTypeClose:
					err := ws.conn.WriteMessage(websocket.CloseMessage,
						websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						remoteAddr := ws.conn.RemoteAddr().String()
						Error("WSClient,WriteMessage error:%s,remote addr:%s,msg type:%d",
							err.Error(), remoteAddr, msg.Type)
					}
					ws.onClose()
					goto QUIT
				}
			case <-keepAliveTimer.C:
				err := ws.conn.WriteMessage(websocket.BinaryMessage, nil)
				if err != nil {
					remoteAddr := ws.conn.RemoteAddr().String()
					Error("WSClient,keepAliveTimer WriteMessage error:%s,remote addr:%s",
						err.Error(), remoteAddr)
					break WriteDone
				}

			case <-reconnectChan:
				break WriteDone
			}
		}
	Reconnect:
		Warning("process reconnecting ...")
		{
			ws.rwLock.RLock()
			isForceQuit := ws.isForceQuit
			ws.rwLock.RUnlock()
			if isForceQuit {
				Warning("WSClient,force quit")
				goto QUIT
			}
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		extraMs := r.Uint32() % 3000 //random add extra reconnect time (0~3s)
		delayTime := ws.reconnectTime + time.Duration(extraMs)*time.Millisecond
		time.Sleep(delayTime)
		ws.disconnectRemote()
		if err := ws.connectRemote(); err != nil {
			Error("WSClient reconnect error:%s", err.Error())
			goto Reconnect
		}
		Info("reconnect success")
	}
QUIT:
	// goto this label,means do not reconnect remote server again.
	Info("WSClient,quit")
}

// RefreshReconnectMsg can refresh the client's reconnect message
func (ws *WSClient) RefreshReconnectMsg(reconnectMsg *WSMessage) {
	ws.rwLock.Lock()
	ws.reconnectMsg = reconnectMsg
	ws.rwLock.Unlock()
}
