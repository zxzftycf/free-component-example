package common

import (
	pb "../grpc"
	"github.com/golang/protobuf/proto"
)

type RedisLockI interface {
	Lock(name string) error
	Unlock(name string)
	MessageLock(name string, extroInfo *pb.MessageExtroInfo, componentName string) error
	MessageUnlock(name string, extroInfo *pb.MessageExtroInfo, componentName string)
}

type LogI interface {
	Info(a ...interface{})
	Error(a ...interface{})
	Debug(a ...interface{})
}

type RouteI interface {
	Call(componentName string, methodName string, request proto.Message, reply proto.Message, extroInfo *pb.MessageExtroInfo) error
	CallAnyReply(componentName string, methodName string, request proto.Message, extroInfo *pb.MessageExtroInfo) (proto.Message, error)
	CallByIp(ip string, componentName string, methodName string, request proto.Message, reply proto.Message, extroInfo *pb.MessageExtroInfo) error
}

type PushI interface {
	Push(request proto.Message, uid string)
	Broadcast(request proto.Message)
	Kick(uid string)
	KickByConnId(connId string, ip string)
	OnLine(userId string, clientConnId string, clientConnIp string) error
	OffLine(uid string)
}
