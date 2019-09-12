package base

import (
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"

	"../common"
	pb "../grpc"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

func init() {
	common.AllComponentMap["Http"] = &Http{}
}

type Http struct {
	Base
}

func (self *Http) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)
	go func() {
		http.HandleFunc("/http/api", self.httpHandler)
		http.ListenAndServe(":"+(*self.Config)["port"], nil)
	}()

	return
}

func (self *Http) httpHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			common.LogError("Http httpHandler has err", err)
			http.Error(w, "", http.StatusInternalServerError)
			debug.PrintStack()
		}
	}()

	messageStartTime := time.Now().UnixNano() / 1e6

	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	//w.Header().Set("Content-type", "text/plain; charset=utf-8")
	//w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, PUT, POST, DELETE")
	ioR, err := ioutil.ReadAll(r.Body)
	//common.LogDebug("Http httpHandler ioR 1", ioR, r.Header, r.Body)
	if err != nil {
		common.LogError("Http httpHandler ReadAll has err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	//common.LogDebug("Http httpHandler ioR", ioR)
	clientRequest := &pb.ClientRequest{}
	err = proto.Unmarshal(ioR, clientRequest)
	if err != nil {
		common.LogError("Http httpHandler proto.Unmarshal has err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	common.LogDebug("Http httpHandler clientRequest", clientRequest)

	componentName := clientRequest.GetComponentName()
	methodName := clientRequest.GetMethodName()
	if common.ClientInterfaceMap[componentName+"."+methodName] != true {
		common.LogError("Http: httpHandler method not found", componentName, methodName)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	clientMessage, err := ptypes.Empty(clientRequest.GetMessageContent())
	if err != nil {
		common.LogError("Http httpHandler ptypes.Empty has err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	err = ptypes.UnmarshalAny(clientRequest.GetMessageContent(), clientMessage)
	if err != nil {
		common.LogError("Http httpHandler ptypes.UnmarshalAny has err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	messageExtroInfo := &pb.MessageExtroInfo{}
	//messageExtroInfo.ClientConnId = c.id
	//messageExtroInfo.ClientConnIp = serverIp + ":" + grpcPort
	//messageExtroInfo.UserId = c.userId
	reply, err := common.Router.CallAnyReply(componentName, methodName, clientMessage, messageExtroInfo)
	if err != nil {
		common.LogError("Http httpHandler common.Router.Call has err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	replyMessageName := proto.MessageName(reply)
	clientReplyContent, err := ptypes.MarshalAny(reply)
	if err != nil {
		common.LogError("Http httpHandler MarshalAny has err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	clientReply := &pb.ClientReply{}
	clientReply.MessageName = replyMessageName
	clientReply.MessageContent = clientReplyContent
	byteInfo, err := proto.Marshal(clientReply)
	if err != nil {
		common.LogError("Http httpHandler Marshal has err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(byteInfo)

	messageEndTime := time.Now().UnixNano() / 1e6
	costTime := messageEndTime - messageStartTime
	common.LogDebug("Http httpHandler message handle ok", costTime, componentName, methodName)
}
