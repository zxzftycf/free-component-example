package base

import (
	"reflect"
	"runtime/debug"

	"../common"
	pb "../grpc"
	"github.com/golang/protobuf/proto"
)

func init() {
	common.AllComponentMap["Route"] = &Route{}
}

type Route struct {
	common.RouteI
	Base
}

func (self *Route) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)

	return
}

func (self *Route) CallAnyReply(componentName string, methodName string, request proto.Message, extroInfo *pb.MessageExtroInfo) (proto.Message, error) {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("Route CallAnyReply has err", err)
			debug.PrintStack()
		}
	}()
	reply, err := self.realCall(componentName, methodName, request, extroInfo)
	return reply, err
}

func (self *Route) Call(componentName string, methodName string, request proto.Message, reply proto.Message, extroInfo *pb.MessageExtroInfo) error {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("Route Call has err", err)
			debug.PrintStack()
		}
	}()
	replyTemp, err := self.realCall(componentName, methodName, request, extroInfo)
	if err != nil {
		return err
	}
	reply.Reset()
	proto.Merge(reply, replyTemp)
	return nil
}

func (self *Route) realCall(componentName string, methodName string, request proto.Message, extroInfo *pb.MessageExtroInfo) (proto.Message, error) {
	var reply proto.Message
	if common.ComponentMap[componentName] != nil {
		methodArgs := []reflect.Value{reflect.ValueOf(request), reflect.ValueOf(extroInfo)}
		rst := reflect.ValueOf(common.ComponentMap[componentName]).MethodByName(methodName).Call(methodArgs)
		if rst[0].Interface() != nil {
			reply = rst[0].Interface().(proto.Message)
		}
		var err error
		if rst[1].Interface() != nil {
			err = rst[1].Interface().(error)
		} else {
			err = nil
		}
		return reply, err
	}

	grpcComponentInterface := common.ComponentMap["GRPC"]
	grpcComponent, _ := grpcComponentInterface.(*GRPC)
	reply, err := grpcComponent.SendMessage(componentName, methodName, request, extroInfo)
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (self *Route) CallByIp(ip string, componentName string, methodName string, request proto.Message, reply proto.Message, extroInfo *pb.MessageExtroInfo) error {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("Route CallByIp has err", err)
			debug.PrintStack()
		}
	}()
	replyTemp, err := self.realCall(componentName, methodName, request, extroInfo)
	if err != nil {
		return err
	}
	reply.Reset()
	proto.Merge(reply, replyTemp)
	return nil
}

func (self *Route) realCallByIp(ip string, componentName string, methodName string, request proto.Message, extroInfo *pb.MessageExtroInfo) (proto.Message, error) {
	var reply proto.Message
	grpcComponentInterface := common.ComponentMap["GRPC"]
	grpcComponent, _ := grpcComponentInterface.(*GRPC)
	reply, err := grpcComponent.SendMessageByIp(ip, componentName, methodName, request, extroInfo)
	if err != nil {
		return reply, err
	}
	return reply, nil
}
