package logic

import (
	"../base"
	"../common"
	pb "../grpc"
	"github.com/golang/protobuf/proto"
)

func init() {
	common.AllComponentMap["PlayerInfo"] = &PlayerInfo{}
}

type PlayerInfo struct {
	Base base.Base
}

func (self *PlayerInfo) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)

	return
}

func (self *PlayerInfo) LoadPlayer(request *pb.LoadPlayerRequest, extroInfo *pb.MessageExtroInfo) (*pb.LoadPlayerReply, error) {
	uuid := request.GetUuid()

	common.Locker.MessageLock(common.Message_Lock_Player+uuid, extroInfo, self.Base.ComponentName)
	defer common.Locker.MessageUnlock(common.Message_Lock_Player+uuid, extroInfo, self.Base.ComponentName)

	common.LogDebug("PlayerInfo LoadPlayer extroInfo", extroInfo)
	common.LogDebug("PlayerInfo LoadPlayer request", request)

	reply := &pb.LoadPlayerReply{}

	//先看redis中有没有
	redisRequest := &pb.RedisMessage{}
	redisRequest.Type = pb.RedisMessageType_GetByte
	redisRequest.Table = common.Redis_PlayerInfo_Table
	redisRequest.Key = uuid
	redisReply := &pb.RedisMessage{}
	err := common.Router.Call("Redis", "Get", redisRequest, redisReply, extroInfo)
	if err != nil {
		common.LogError("PlayerInfo LoadPlayer Call Redis Get has err", err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	playerInfoByte := redisReply.ValueByte
	if playerInfoByte != nil {
		playerInfo := &pb.PlayerInfo{}
		err = proto.Unmarshal(playerInfoByte, playerInfo)
		if err != nil {
			common.LogError("PlayerInfo LoadPlayer proto.Unmarshal has err", err)
			reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
			return reply, nil
		}
		reply.PlayerInfo = playerInfo
		common.LogInfo("PlayerInfo LoadPlayer from redis ok:", uuid)
		return reply, nil
	}
	//没有就去数据库里读
	mysqlKVRequest := &pb.MysqlKVMessage{}
	mysqlKVRequest.TableName = common.Mysql_PlayerInfo_Table
	mysqlKVRequest.Uuid = uuid
	mysqlKVReply := &pb.MysqlKVMessage{}
	err = common.Router.Call("Mysql", "QueryKV", mysqlKVRequest, mysqlKVReply, extroInfo)
	if err != nil {
		common.LogError("PlayerInfo LoadPlayer Call Mysql QueryKV has err", uuid, err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	if mysqlKVReply.GetErrMessage() != nil {
		reply.ErrMessage = mysqlKVReply.GetErrMessage()
		return reply, nil
	}
	playerInfoByte = mysqlKVReply.GetInfo()
	if playerInfoByte == nil {
		common.LogError("PlayerInfo LoadPlayer playerInfoByte == nil in mysql", uuid)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	playerInfo := &pb.PlayerInfo{}
	err = proto.Unmarshal(playerInfoByte, playerInfo)
	if err != nil {
		common.LogError("PlayerInfo LoadPlayer proto.Unmarshal after mysql has err", err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	//有就存到redis里一份
	redisRequest = &pb.RedisMessage{}
	redisRequest.Type = pb.RedisMessageType_SetByte
	redisRequest.Table = common.Redis_PlayerInfo_Table
	redisRequest.Key = uuid
	redisRequest.ValueByte = playerInfoByte
	redisReply = &pb.RedisMessage{}
	err = common.Router.Call("Redis", "Set", redisRequest, redisReply, extroInfo)
	if err != nil {
		common.LogError("PlayerInfo LoadPlayer Call Redis Set has err", err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	reply.PlayerInfo = playerInfo
	common.LogInfo("PlayerInfo LoadPlayer from mysql ok:", uuid)
	return reply, nil
}

func (self *PlayerInfo) NewPlayer(request *pb.NewPlayerRequest, extroInfo *pb.MessageExtroInfo) (*pb.NewPlayerReply, error) {
	uuid := request.GetUuid()
	shortId := request.GetShortId()

	common.Locker.MessageLock(common.Message_Lock_Player+uuid, extroInfo, self.Base.ComponentName)
	defer common.Locker.MessageUnlock(common.Message_Lock_Player+uuid, extroInfo, self.Base.ComponentName)

	common.LogDebug("PlayerInfo NewPlayer extroInfo", extroInfo)
	common.LogDebug("PlayerInfo NewPlayer request", request)

	reply := &pb.NewPlayerReply{}

	playerInfo := &pb.PlayerInfo{}
	playerInfo.Name = ""
	playerInfo.Balance = 100
	playerInfo.Uuid = ""
	playerInfo.ShortId = ""

	playerInfoByte, err := proto.Marshal(playerInfo)
	if err != nil {
		common.LogError("PlayerInfo: NewPlayer Marshal has err", uuid, shortId, err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	//先存数据库
	mysqlKVRequest := &pb.MysqlKVMessage{}
	mysqlKVRequest.TableName = common.Mysql_PlayerInfo_Table
	mysqlKVRequest.Uuid = uuid
	mysqlKVRequest.ShortId = shortId
	mysqlKVRequest.Info = playerInfoByte
	mysqlKVReply := &pb.MysqlKVMessage{}
	err = common.Router.Call("Mysql", "InsertKV", mysqlKVRequest, mysqlKVReply, extroInfo)
	if err != nil {
		common.LogError("PlayerInfo NewPlayer Call Mysql InsertKV has err", uuid, shortId, err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}
	if mysqlKVReply.GetErrMessage() != nil {
		reply.ErrMessage = mysqlKVReply.GetErrMessage()
		return reply, nil
	}
	//再存redis
	redisRequest := &pb.RedisMessage{}
	redisRequest.Type = pb.RedisMessageType_SetByte
	redisRequest.Table = common.Redis_PlayerInfo_Table
	redisRequest.Key = uuid
	redisRequest.ValueByte = playerInfoByte
	redisReply := &pb.RedisMessage{}
	err = common.Router.Call("Redis", "Set", redisRequest, redisReply, extroInfo)
	if err != nil {
		common.LogError("PlayerInfo NewPlayer Call Redis Set has err", err)
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_ServerError, "")
		return reply, nil
	}

	reply.PlayerInfo = playerInfo
	common.LogInfo("PlayerInfo NewPlayer ok:", uuid, shortId)
	return reply, nil
}
