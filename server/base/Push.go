package base

import (
	"../common"
	pb "../grpc"
	"github.com/golang/protobuf/proto"
)

func init() {
	common.AllComponentMap["Push"] = &Push{}
}

type Push struct {
	Base
	common.PushI
}

func (self *Push) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)

	return
}

func (self *Push) OnLine(userId string, clientConnId string, clientConnIp string) error {
	request := &pb.EmptyMessage{}
	reply := &pb.EmptyMessage{}
	extroInfo := &pb.MessageExtroInfo{}
	extroInfo.UserId = userId
	extroInfo.ClientConnId = clientConnId
	err := common.Router.CallByIp(clientConnIp, "WebSocket", "UserLoginOk", request, reply, extroInfo)
	if err != nil {
		common.LogError("Push OnLine CallByIp WebSocket UserLoginOk has err", err)
		return err
	}
	request1 := &pb.RedisMessage{}
	request1.Type = pb.RedisMessageType_SetString
	request1.Table = common.Redis_OnlineUser_Table
	request1.Key = userId
	request1.ValueString = clientConnIp
	reply1 := &pb.RedisMessage{}
	err = common.Router.Call("Redis", "Set", request1, reply1, extroInfo)
	if err != nil {
		common.LogError("Push OnLine Call Redis Set has err", err)
		return err
	}
	return nil
}

func (self *Push) OffLine(uid string) {
	request1 := &pb.RedisMessage{}
	request1.Table = common.Redis_OnlineUser_Table
	request1.Key = uid
	reply1 := &pb.RedisMessage{}
	extroInfo := &pb.MessageExtroInfo{}
	err := common.Router.Call("Redis", "Delete", request1, reply1, extroInfo)
	if err != nil {
		common.LogError("Push OnLine Call Redis Delete has err", err)
		return
	}
}

func (self *Push) Push(pushMessage proto.Message, uid string) {
	request := &pb.RedisMessage{}
	request.Type = pb.RedisMessageType_GetString
	request.Table = common.Redis_OnlineUser_Table
	request.Key = uid
	reply := &pb.RedisMessage{}
	extroInfo := &pb.MessageExtroInfo{}
	err := common.Router.Call("Redis", "Get", request, reply, extroInfo)
	if err != nil {
		common.LogError("Push Push Call Redis Get has err", err)
		return
	}
	ip := reply.ValueString
	if ip == "" {
		common.LogInfo("Push Push user not login", uid)
		return
	}
	extroInfo.UserId = uid
	reply1 := &pb.EmptyMessage{}
	err = common.Router.CallByIp(ip, "WebSocket", "Push", pushMessage, reply1, extroInfo)
	if err != nil {
		common.LogError("Push Push CallByIp WebSocket Push has err", err)
		return
	}
}

func (self *Push) Broadcast(request proto.Message) {
	findComponentInterface := common.ComponentMap["Find"]
	findComponent, _ := findComponentInterface.(*Find)
	ips, err := findComponent.FindAllComponent("WebSocket")
	if err != nil {
		common.LogError("Push Broadcast FindAllComponent has err", err)
		return
	}
	if len(ips) <= 0 {
		common.LogInfo("Push Broadcast no websocket component")
		return
	}

	for ip, _ := range ips {
		go func(websocketIp string) {
			reply := &pb.EmptyMessage{}
			extroInfo := &pb.MessageExtroInfo{}
			err = common.Router.CallByIp(websocketIp, "WebSocket", "Broadcast", request, reply, extroInfo)
			if err != nil {
				common.LogError("Push Push CallByIp WebSocket Push has err", err)
				return
			}
		}(ip)
	}
}

func (self *Push) Kick(uid string) {
	request := &pb.RedisMessage{}
	request.Type = pb.RedisMessageType_GetString
	request.Table = common.Redis_OnlineUser_Table
	request.Key = uid
	reply := &pb.RedisMessage{}
	extroInfo := &pb.MessageExtroInfo{}
	err := common.Router.Call("Redis", "Get", request, reply, extroInfo)
	if err != nil {
		common.LogError("Push Kick Call Redis Get has err", err)
		return
	}
	ip := reply.ValueString
	if ip == "" {
		common.LogInfo("Push Kick user not login", uid)
		return
	}
	extroInfo.UserId = uid
	reply1 := &pb.EmptyMessage{}
	err = common.Router.CallByIp(ip, "WebSocket", "Kick", request, reply1, extroInfo)
	if err != nil {
		common.LogError("Push Kick CallByIp WebSocket Push has err", err)
		return
	}
}

func (self *Push) KickByConnId(connId string, ip string) {
	extroInfo := &pb.MessageExtroInfo{}
	extroInfo.ClientConnId = connId
	reply1 := &pb.EmptyMessage{}
	err := common.Router.CallByIp(ip, "WebSocket", "Kick", reply1, reply1, extroInfo)
	if err != nil {
		common.LogError("Push KickByConnId CallByIp WebSocket Push has err", err)
		return
	}
}
