package base

import (
	"errors"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/golang/protobuf/proto"

	"../common"
	pb "../grpc"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func init() {
	common.AllComponentMap["WebSocket"] = &WebSocket{}
}

type ClientManager struct {
	clients  map[string]*Client
	lManager sync.Mutex
}

type Client struct {
	id        string
	userId    string
	socket    *websocket.Conn
	send      chan []byte
	lConn     sync.Mutex
	startTime int64
}

type WebSocket struct {
	Base
	manager ClientManager
}

var serverIp string
var grpcPort string

func (self *WebSocket) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)
	self.manager = ClientManager{
		clients: make(map[string]*Client),
	}

	conn, err := net.Dial("udp", "www.google.com.hk:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	serverIp = strings.Split(conn.LocalAddr().String(), ":")[0]
	grpcPort = (*self.Config)["grpc_port"]

	go func() {
		http.HandleFunc("/ws", self.wsPage)
		http.ListenAndServe((*self.Config)["listen_url"], nil)

	}()
	common.LogInfo("WebSocketComponent listen on :", (*self.Config)["listen_url"])
	return
}

func (m *ClientManager) updateClient(clientId string, userId string) (*Client, error) {
	m.lManager.Lock()
	defer m.lManager.Unlock()

	if _, ok := m.clients[userId]; ok {
		return nil, errors.New("WebSocket updateClient userId already in clients:" + userId)
	}
	if _, ok := m.clients[clientId]; !ok {
		return nil, errors.New("WebSocket updateClient clientId err," + clientId)
	}
	m.clients[userId] = m.clients[clientId]
	delete(m.clients, clientId)
	return m.clients[userId], nil
}

func (m *ClientManager) validClient(c *Client) bool {
	m.lManager.Lock()
	defer m.lManager.Unlock()

	if _, ok := m.clients[c.id]; ok {
		return true
	}
	if _, ok := m.clients[c.userId]; ok {
		return true
	}
	return false
}

func (m *ClientManager) registerClient(c *Client) {
	m.lManager.Lock()
	defer m.lManager.Unlock()

	m.clients[c.id] = c
	common.LogInfo("WebSocketComponent: registerClient", c.id, c.userId)
}

func (m *ClientManager) unregisterClient(c *Client) {
	m.lManager.Lock()
	defer m.lManager.Unlock()

	delete(m.clients, c.id)
	delete(m.clients, c.userId)

	common.LogInfo("WebSocketComponent: unregisterClient", c.id, c.userId)
}

func (c *Client) update(userId string) {
	c.lConn.Lock()
	defer c.lConn.Unlock()

	c.userId = userId
}

func (c *Client) close() {
	c.lConn.Lock()
	defer c.lConn.Unlock()

	common.LogInfo("WebSocketComponent: Client close", c.id, c.userId)
	c.socket.Close()
	close(c.send)
	c.userId = ""
}

func (c *Client) sendMessage(messageName string, message *any.Any) {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("WebSocket sendMessage has err", err)
			debug.PrintStack()
		}
	}()

	c.lConn.Lock()
	defer c.lConn.Unlock()

	/*if c.userId == "" {
		return
	}*/

	clientReply := &pb.ClientReply{}
	clientReply.MessageName = messageName
	clientReply.MessageContent = message
	byteInfo, err := proto.Marshal(clientReply)
	if err != nil {
		common.LogError("WebSocketComponent: sendMessage has err", c.id, c.userId, err)
		return
	}

	c.send <- byteInfo
}

func (c *Client) receiveMessage(message []byte) error {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("WebSocket receiveMessage has err", err)
			debug.PrintStack()
		}
	}()

	c.lConn.Lock()
	defer c.lConn.Unlock()

	messageStartTime := time.Now().UnixNano() / 1e6

	clientRequest := &pb.ClientRequest{}
	err := proto.Unmarshal(message, clientRequest)
	if err != nil {
		return err
	}
	componentName := clientRequest.GetComponentName()
	methodName := clientRequest.GetMethodName()
	if common.ClientInterfaceMap[componentName+"."+methodName] != true {
		common.LogError("WebSocketComponent: receiveMessage method not found", c.id, c.userId, componentName, methodName)
		return errors.New("method not found")
	}

	clientMessage, err := ptypes.Empty(clientRequest.GetMessageContent())
	if err != nil {
		common.LogError("WebSocketComponent ptypes.Empty has err", c.id, c.userId, err)
		return err
	}
	err = ptypes.UnmarshalAny(clientRequest.GetMessageContent(), clientMessage)
	if err != nil {
		common.LogError("WebSocketComponent ptypes.UnmarshalAny has err", c.id, c.userId, err)
		return err
	}

	messageExtroInfo := &pb.MessageExtroInfo{}
	messageExtroInfo.ClientConnId = c.id
	messageExtroInfo.ClientConnIp = serverIp + ":" + grpcPort
	messageExtroInfo.UserId = c.userId
	reply, err := common.Router.CallAnyReply(componentName, methodName, clientMessage, messageExtroInfo)
	if err != nil {
		common.LogError("WebSocketComponent common.Router.Call has err", c.id, c.userId, err)
		return nil
	}
	replyMessageName := proto.MessageName(reply)
	clientReplyContent, err := ptypes.MarshalAny(reply)
	if err != nil {
		common.LogError("WebSocketComponent MarshalAny has err", c.id, c.userId, err)
		return err
	}
	go c.sendMessage(replyMessageName, clientReplyContent)

	messageEndTime := time.Now().UnixNano() / 1e6
	costTime := messageEndTime - messageStartTime
	common.LogDebug("message handle ok", costTime, componentName, methodName)

	return nil
}

func (self *WebSocket) UserLoginOk(message *pb.EmptyMessage, extroInfo *pb.MessageExtroInfo) (*pb.EmptyMessage, error) {
	userId := extroInfo.GetUserId()
	clientConnId := extroInfo.GetClientConnId()
	reply := &pb.EmptyMessage{}
	c, err := self.manager.updateClient(clientConnId, userId)
	if err != nil {
		common.LogError("WebSocket UserLoginOk updateClient has err", err)
		return reply, err
	}
	c.update(userId)
	return reply, nil
}

func (self *WebSocket) Push(message proto.Message, extroInfo *pb.MessageExtroInfo) (*pb.EmptyMessage, error) {
	reply := &pb.EmptyMessage{}
	userId := extroInfo.GetUserId()
	if _, ok := self.manager.clients[userId]; !ok {
		return reply, nil
	}

	replyMessageName := proto.MessageName(message)
	clientReplyContent, err := ptypes.MarshalAny(message)
	if err != nil {
		common.LogError("WebSocket Push MarshalAny has err", err)
		return nil, err
	}
	go self.manager.clients[userId].sendMessage(replyMessageName, clientReplyContent)
	return reply, nil
}

func (self *WebSocket) Kick(message proto.Message, extroInfo *pb.MessageExtroInfo) (*pb.EmptyMessage, error) {
	reply := &pb.EmptyMessage{}
	userId := extroInfo.GetUserId()
	clientConnId := extroInfo.GetClientConnId()
	if userId != "" {
		if _, ok := self.manager.clients[userId]; !ok {
			common.LogInfo("WebSocket Kick user not online:", userId, clientConnId)
			return reply, nil
		}
		go self.clientClose(self.manager.clients[userId])
		return reply, nil
	}
	if clientConnId != "" {
		if _, ok := self.manager.clients[clientConnId]; !ok {
			common.LogInfo("WebSocket Kick clientConn not online:", userId, clientConnId)
			return reply, nil
		}
		go self.clientClose(self.manager.clients[clientConnId])
		return reply, nil
	}
	return reply, errors.New("WebSocket Kick userId or clientConnId is nil")
}

func (self *WebSocket) Broadcast(message proto.Message, extroInfo *pb.MessageExtroInfo) (*pb.EmptyMessage, error) {
	reply := &pb.EmptyMessage{}
	for _, conn := range self.manager.clients {
		replyMessageName := proto.MessageName(message)
		clientReplyContent, err := ptypes.MarshalAny(message)
		if err != nil {
			common.LogInfo("WebSocket Broadcast MarshalAny has err", err)
			continue
		}
		go conn.sendMessage(replyMessageName, clientReplyContent)
	}
	return reply, nil
}

func (self *WebSocket) clientClose(c *Client) {
	isValid := self.manager.validClient(c)
	if isValid == false {
		return
	}
	if c.userId != "" {
		common.Pusher.OffLine(c.userId)
	}
	self.manager.unregisterClient(c)
	c.close()
}

func (self *WebSocket) clientRead(c *Client) {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("WebSocket clientRead has err", err)
			debug.PrintStack()
		}
		self.clientClose(c)
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			//common.LogError("WebSocketComponent clientRead ReadMessage has err", err)
			break
		}
		err = c.receiveMessage(message)
		if err != nil {
			common.LogError("WebSocketComponent clientRead receiveMessage has err", err)
			break
		}
	}
}

func (self *WebSocket) clientWrite(c *Client) {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("WebSocket clientWrite has err", err)
			debug.PrintStack()
		}
		self.clientClose(c)
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				//common.LogError("WebSocketComponent clientWrite c.send has err")
				return
			}
			c.socket.WriteMessage(websocket.BinaryMessage, message)
		}
	}
}

func (self *WebSocket) wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		common.LogError("WebSocketComponent wsPage conn has err", err)
		return
	}
	connUuid, err := uuid.NewV4()
	if err != nil {
		http.NotFound(res, req)
		common.LogError("WebSocketComponent wsPage connUuid has err", err)
		return
	}
	client := &Client{id: connUuid.String(), userId: "", socket: conn, send: make(chan []byte), startTime: time.Now().Unix()}

	common.LogInfo("WebSocketComponent: Client open", client.id, client.userId)

	self.manager.registerClient(client)

	go self.clientWrite(client)
	go self.clientRead(client)

	//10秒内没有登陆成功就关闭这个链接
	tm := time.NewTimer(time.Second * 10)
	go func() {
		<-tm.C
		tm.Stop()
		if client.userId == "" {
			common.LogInfo("WebSocket wsPage one socket invalid:", client.id)
			self.clientClose(client)
		}
	}()
}
