package logic

import (
	"../base"
	"../common"
	pb "../grpc"
)

func init() {
	common.AllComponentMap["Register"] = &Register{}
}

type Register struct {
	Base base.Base
}

func (self *Register) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)

	return
}

func (self *Register) Register(request *pb.RegisterRequest, extroInfo *pb.MessageExtroInfo) (*pb.RegisterReply, error) {
	account := request.GetAccount()
	password := request.GetPassword()
	common.Locker.MessageLock(common.Message_Lock_Account_Register+account, extroInfo, self.Base.ComponentName)
	defer common.Locker.MessageUnlock(common.Message_Lock_Account_Register+account, extroInfo, self.Base.ComponentName)

	common.LogDebug("Register Register extroInfo", extroInfo)
	common.LogDebug("Register Register request", request)
	reply := &pb.RegisterReply{}
	mysqlRequest := &pb.MysqlAccountInfo{}
	mysqlRequest.Account = account
	mysqlRequest.Password = password
	mysqlReply := &pb.MysqlAccountInfo{}
	err := common.Router.Call("Mysql", "NewAccount", mysqlRequest, mysqlReply, extroInfo)
	if err != nil {
		common.LogError("Register Register Call Mysql NewAccount has err", account, err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	if mysqlReply.GetErrMessage() != nil {
		reply.ErrMessage = mysqlReply.GetErrMessage()
		return reply, nil
	}

	common.LogDebug("Register Register ok", mysqlReply)
	return reply, nil
}
