package base

import (
	"context"
	"errors"
	"net"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"../common"
	pb "../grpc"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

func init() {
	common.AllComponentMap["GRPC"] = &GRPC{}
}

type ComponentConn struct {
	Conn          *grpc.ClientConn
	ComponentName string
	Ip            string
}

type ComponentConns map[string]ComponentConn

type GRPC struct {
	Base
	ConnPool          map[string]ComponentConns
	lConnPool         sync.Mutex
	refreshConnsTimer *time.Ticker
}

func (self *GRPC) HandleMessage(ctx context.Context, in *pb.HandleMessageRequest) (*pb.HandleMessageReply, error) {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("GRPC HandleMessage has err", err)
			debug.PrintStack()
		}
	}()
	grpcReply := &pb.HandleMessageReply{}
	componentName := in.GetComponentName()
	methodName := in.GetMethodName()
	inMessage := in.GetMessageContent()
	extroInfo := in.GetExtroInfo()
	realInMessage, err := ptypes.Empty(inMessage)
	if err != nil {
		common.LogError("GRPC HandleMessage ptypes.Empty has err", err)
		return grpcReply, err
	}
	err = ptypes.UnmarshalAny(inMessage, realInMessage)
	if err != nil {
		common.LogError("GRPC HandleMessage ptypes.UnmarshalAny has err", err)
		return grpcReply, err
	}
	if common.ComponentMap[componentName] != nil {
		methodArgs := []reflect.Value{reflect.ValueOf(realInMessage), reflect.ValueOf(extroInfo)}
		rst := reflect.ValueOf(common.ComponentMap[componentName]).MethodByName(methodName).Call(methodArgs)
		var reply proto.Message
		if rst[0].Interface() != nil {
			reply = rst[0].Interface().(proto.Message)
		} else {
			reply = nil
		}
		var err error
		if rst[1].Interface() != nil {
			err = rst[1].Interface().(error)
		} else {
			err = nil
		}
		if err != nil {
			common.LogError("GRPC HandleMessage Call has err", err)
			return grpcReply, err
		}

		replyAny, err := ptypes.MarshalAny(reply)
		if err != nil {
			common.LogError("GRPC HandleMessage MarshalAny(reply) has err", err)
			return grpcReply, err
		}
		grpcReply.MessageContent = replyAny
		return grpcReply, nil
	}

	return grpcReply, errors.New("GRPC HandleMessage componentName err:" + componentName)
}

func (self *GRPC) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)
	self.ConnPool = make(map[string]ComponentConns)
	go func() {
		listen, err := net.Listen("tcp", (*self.Config)["listen_url"])
		if err != nil {
			panic(err)
		}
		//实现gRPC Server
		s := grpc.NewServer()
		//注册helloServer为客户端提供服务
		pb.RegisterGRPCComponentServer(s, self) //内部调用了s.RegisterServer()
		common.LogInfo("GRPCComponent Listen on", (*self.Config)["listen_url"])

		s.Serve(listen)
	}()

	self.refreshConnsTimer = time.NewTicker(5 * time.Second)
	go func(t *time.Ticker) {
		for {
			<-t.C
			self.RefreshConn()
		}
	}(self.refreshConnsTimer)

	return
}

func (self *GRPC) SendMessage(componentName string, methodName string, request proto.Message, extroInfo *pb.MessageExtroInfo) (proto.Message, error) {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("GRPC SendMessage has err", err)
			debug.PrintStack()
		}
	}()
	componentConn, err := self.NewConn(componentName)
	var reply proto.Message
	if err != nil {
		common.LogError("GRPC SendMessage NewConn has err", err)
		return reply, err
	}
	grpcContent, err := ptypes.MarshalAny(request)
	if err != nil {
		common.LogError("GRPC SendMessage MarshalAny(request) has err", err)
		return reply, err
	}
	//初始化客户端
	c := pb.NewGRPCComponentClient(componentConn.Conn)
	//调用方法
	reqBody := &pb.HandleMessageRequest{}
	reqBody.ComponentName = componentName
	reqBody.MethodName = methodName
	reqBody.MessageContent = grpcContent
	reqBody.ExtroInfo = extroInfo
	ctx1, cel := context.WithTimeout(context.Background(), time.Second*10)
	defer cel()
	r, err := c.HandleMessage(ctx1, reqBody)
	if err != nil {
		common.LogError("GRPC SendMessage Call HandleMessage has err", componentName, methodName, request, extroInfo, err)
		return reply, err
	}
	grpcReply := r.GetMessageContent()
	realGrpcReply, err := ptypes.Empty(grpcReply)
	if err != nil {
		common.LogError("GRPC SendMessage ptypes.Empty has err", err)
		return grpcReply, err
	}
	err = ptypes.UnmarshalAny(grpcReply, realGrpcReply)
	if err != nil {
		common.LogError("GRPC SendMessage UnmarshalAny(grpcReply, reply) has err", err)
		return reply, err
	}
	return realGrpcReply, nil
}

func (self *GRPC) SendMessageByIp(ip string, componentName string, methodName string, request proto.Message, extroInfo *pb.MessageExtroInfo) (proto.Message, error) {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("GRPC SendMessageByIp has err", err)
			debug.PrintStack()
		}
	}()
	componentConn, err := self.NewConnByIp(ip, componentName)
	var reply proto.Message
	if err != nil {
		common.LogError("GRPC SendMessageByIp NewConn has err", err)
		return reply, err
	}
	grpcContent, err := ptypes.MarshalAny(request)
	if err != nil {
		common.LogError("GRPC SendMessageByIp MarshalAny(request) has err", err)
		return reply, err
	}
	//初始化客户端
	c := pb.NewGRPCComponentClient(componentConn.Conn)
	//调用方法
	reqBody := &pb.HandleMessageRequest{}
	reqBody.ComponentName = componentName
	reqBody.MethodName = methodName
	reqBody.MessageContent = grpcContent
	reqBody.ExtroInfo = extroInfo
	ctx1, cel := context.WithTimeout(context.Background(), time.Second*10)
	defer cel()
	r, err := c.HandleMessage(ctx1, reqBody)
	if err != nil {
		common.LogError("GRPC SendMessageByIp Call HandleMessage has err", componentName, methodName, request, extroInfo, err)
		return reply, err
	}
	grpcReply := r.GetMessageContent()
	realGrpcReply, err := ptypes.Empty(grpcReply)
	if err != nil {
		common.LogError("GRPC SendMessageByIp ptypes.Empty has err", err)
		return grpcReply, err
	}
	err = ptypes.UnmarshalAny(grpcReply, realGrpcReply)
	if err != nil {
		common.LogError("GRPC SendMessageByIp UnmarshalAny(grpcReply, reply) has err", err)
		return reply, err
	}
	return realGrpcReply, nil
}

func (self *GRPC) DeleteConn(conn *ComponentConn) {
	self.lConnPool.Lock()
	defer self.lConnPool.Unlock()

	self.deleteConnNoLock(conn)
}

func (self *GRPC) deleteConnNoLock(conn *ComponentConn) {
	if self.ConnPool[conn.ComponentName] == nil {
		return
	}

	conn.Conn.Close()
	delete(self.ConnPool[conn.ComponentName], conn.Ip)
}

func (self *GRPC) NewConn(componentName string) (*ComponentConn, error) {
	self.lConnPool.Lock()
	defer self.lConnPool.Unlock()

	if self.ConnPool[componentName] == nil || len(self.ConnPool[componentName]) <= 0 {
		findComponentInterface := common.ComponentMap["Find"]
		findComponent, _ := findComponentInterface.(*Find)
		ip, err := findComponent.FindComponent(componentName)
		if err != nil {
			return nil, err
		}

		conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return nil, err
		}

		if self.ConnPool[componentName] == nil {
			self.ConnPool[componentName] = make(ComponentConns)
		}

		componentConn := ComponentConn{}
		componentConn.Conn = conn
		componentConn.ComponentName = componentName
		componentConn.Ip = ip
		self.ConnPool[componentName][ip] = componentConn
		return &componentConn, nil
	}

	conns := self.ConnPool[componentName]
	for _, conn := range conns {
		return &conn, nil
	}
	return nil, errors.New("GRPC NewConn no conn:" + componentName)
}

func (self *GRPC) NewConnByIp(ip string, componentName string) (*ComponentConn, error) {
	self.lConnPool.Lock()
	defer self.lConnPool.Unlock()

	findComponentInterface := common.ComponentMap["Find"]
	findComponent, _ := findComponentInterface.(*Find)
	ips, err := findComponent.FindAllComponent(componentName)
	if err != nil {
		return nil, err
	}
	if _, ok := ips[ip]; !ok {
		return nil, errors.New("GRPC NewConnByIp " + componentName + " not has this ip:" + ip)
	}

	if self.ConnPool[componentName] == nil {
		self.ConnPool[componentName] = make(ComponentConns)
	}
	if _, ok := self.ConnPool[componentName][ip]; !ok {
		conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return nil, err
		}

		componentConn := ComponentConn{}
		componentConn.Conn = conn
		componentConn.ComponentName = componentName
		componentConn.Ip = ip
		self.ConnPool[componentName][ip] = componentConn
		return &componentConn, nil
	}

	conn := self.ConnPool[componentName][ip]
	return &conn, nil
}

func (self *GRPC) RefreshConn() {
	self.lConnPool.Lock()
	defer self.lConnPool.Unlock()

	for componentName, _ := range self.ConnPool {
		if self.ConnPool[componentName] == nil {
			self.ConnPool[componentName] = make(ComponentConns)
		}
		findComponentInterface := common.ComponentMap["Find"]
		findComponent, _ := findComponentInterface.(*Find)
		ips, err := findComponent.FindAllComponent(componentName)
		if err != nil {
			continue
		}
		for ip, conn := range self.ConnPool[componentName] {
			if _, ok := ips[ip]; ok {
				continue
			}
			self.deleteConnNoLock(&conn)
		}
		for ip, _ := range ips {
			if _, ok := self.ConnPool[componentName][ip]; ok {
				continue
			}
			conn, err := grpc.Dial(ip, grpc.WithBlock())
			if err != nil {
				continue
			}

			componentConn := ComponentConn{}
			componentConn.Conn = conn
			componentConn.ComponentName = componentName
			componentConn.Ip = ip
			self.ConnPool[componentName][ip] = componentConn
			break
		}
	}

	//common.LogDebug("GRPC RefreshConn ok", self.ConnPool)
}
