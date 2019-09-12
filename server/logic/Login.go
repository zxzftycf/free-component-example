package logic

import (
	"time"

	"../base"
	"../common"
	pb "../grpc"
)

func init() {
	common.AllComponentMap["Login"] = &Login{}
}

type Login struct {
	Base base.Base
}

func (self *Login) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)

	return
}

func (self *Login) Login(request *pb.LoginRequest, extroInfo *pb.MessageExtroInfo) (*pb.LoginReply, error) {
	account := request.GetAccount()
	password := request.GetPassword()

	common.Locker.MessageLock(common.Message_Lock_Account_Login+account, extroInfo, self.Base.ComponentName)
	defer common.Locker.MessageUnlock(common.Message_Lock_Account_Login+account, extroInfo, self.Base.ComponentName)

	common.LogDebug("Login Login extroInfo", extroInfo)
	common.LogDebug("Login Login request", request)

	reply := &pb.LoginReply{}
	mysqlRequest := &pb.MysqlAccountInfo{}
	mysqlRequest.Account = account
	mysqlRequest.Password = password
	mysqlReply := &pb.MysqlAccountInfo{}
	err := common.Router.Call("Mysql", "VerifyAccount", mysqlRequest, mysqlReply, extroInfo)
	if err != nil {
		common.LogError("Login Login Call Mysql VerifyAccount has err", account, err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	if mysqlReply.GetErrMessage() != nil {
		reply.ErrMessage = mysqlReply.GetErrMessage()
		return reply, nil
	}

	common.LogDebug("Login Login VerifyAccount ok")

	uuid := mysqlReply.GetUuid()
	shortId := mysqlReply.GetShortId()
	common.Locker.MessageLock(common.Message_Lock_Player+uuid, extroInfo, self.Base.ComponentName)
	defer common.Locker.MessageUnlock(common.Message_Lock_Player+uuid, extroInfo, self.Base.ComponentName)

	loadPlayerRequest := &pb.LoadPlayerRequest{}
	loadPlayerRequest.Uuid = uuid
	loadPlayerReply := &pb.LoadPlayerReply{}
	err = common.Router.Call("PlayerInfo", "LoadPlayer", loadPlayerRequest, loadPlayerReply, extroInfo)
	if err != nil {
		common.LogError("Login Login Call PlayerInfo LoadPlayer has err", uuid, shortId, err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}

	if loadPlayerReply.GetErrMessage() != nil && loadPlayerReply.GetErrMessage().GetCode() == pb.ErrorCode_DataNotFound {
		reply.ErrMessage = mysqlReply.GetErrMessage()
		return reply, nil
	}

	if loadPlayerReply.GetErrMessage() != nil {
		reply.ErrMessage = mysqlReply.GetErrMessage()
		return reply, nil
	}

	playerInfo := &pb.PlayerInfo{}
	playerInfo.Uuid = "1"

	reply.PlayerInfo = playerInfo
	extroInfo.UserId = "1"
	userId := extroInfo.GetUserId()
	clientConnIp := extroInfo.GetClientConnIp()
	clientConnId := extroInfo.GetClientConnId()
	go func() {
		err := common.Pusher.OnLine(userId, clientConnId, clientConnIp)
		if err != nil {
			common.LogError("Login Login OnLine has err", err)
			go common.Pusher.KickByConnId(clientConnId, clientConnIp)
		}

		time.Sleep(time.Duration(2) * time.Second)
		go common.Pusher.Push(reply, "1")
		go common.Pusher.Push(reply, "2")
		go common.Pusher.Broadcast(reply)
		time.Sleep(time.Duration(2) * time.Second)
		go common.Pusher.Kick("1")
		go common.Pusher.Kick("2")
	}()
	return reply, nil
}
