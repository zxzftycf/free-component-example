/*
	用于redis操作
*/

package base

import (
	"errors"
	"time"

	"../common"
	pb "../grpc"
	"github.com/gomodule/redigo/redis"
)

func init() {
	common.AllComponentMap["Redis"] = &Redis{}
}

type Redis struct {
	Base
	Pool redis.Pool
}

func (self *Redis) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)
	self.Pool = redis.Pool{
		MaxIdle:     16,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", (*self.Config)["redis_host"])
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return
}

func (self *Redis) Set(request *pb.RedisMessage, extroInfo *pb.MessageExtroInfo) (*pb.RedisMessage, error) {
	common.LogDebug("Redis Set request", request)
	conn := self.Pool.Get()
	defer conn.Close()
	table := request.GetTable()
	key := request.GetKey()
	redisType := request.GetType()
	var value interface{}
	reply := &pb.RedisMessage{}
	if redisType == pb.RedisMessageType_SetString {
		value = request.GetValueString()
	} else if redisType == pb.RedisMessageType_SetByte {
		value = request.GetValueByte()
	} else {
		common.LogError("Redis Set message type has err", redisType)
		return reply, errors.New("Redis Set wrong message type")
	}

	realKey := table + ":" + key
	_, err := conn.Do("SET", realKey, value)
	if err != nil {
		common.LogError("Redis Set has err", err)
		return reply, err
	}
	return reply, nil
}

func (self *Redis) Get(request *pb.RedisMessage, extroInfo *pb.MessageExtroInfo) (*pb.RedisMessage, error) {
	common.LogDebug("Redis Get request", request)
	conn := self.Pool.Get()
	defer conn.Close()
	table := request.GetTable()
	key := request.GetKey()
	redisType := request.GetType()
	reply := &pb.RedisMessage{}
	realKey := table + ":" + key
	if redisType == pb.RedisMessageType_GetString {
		reply.ValueString = ""
		res, err := redis.String(conn.Do("GET", realKey))
		if err == redis.ErrNil {
			return reply, nil
		}
		reply.ValueString = res
		if err != nil {
			common.LogError("Redis Get has err", err)
			return reply, err
		}
	} else if redisType == pb.RedisMessageType_GetByte {
		reply.ValueByte = nil
		res, err := redis.Bytes(conn.Do("GET", realKey))
		if err == redis.ErrNil {
			return reply, nil
		}
		if err != nil {
			common.LogError("Redis Get has err", err)
			return reply, err
		}
		reply.ValueByte = res
	} else {
		common.LogError("Redis Get message type has err", redisType)
		return reply, errors.New("Redis Get wrong message type")
	}
	return reply, nil
}

func (self *Redis) Delete(request *pb.RedisMessage, extroInfo *pb.MessageExtroInfo) (*pb.RedisMessage, error) {
	common.LogDebug("Redis Delete request", request)
	conn := self.Pool.Get()
	defer conn.Close()
	table := request.GetTable()
	key := request.GetKey()
	reply := &pb.RedisMessage{}
	realKey := table + ":" + key
	_, err := conn.Do("DEL", realKey)
	if err != nil {
		common.LogError("Redis Delete has err", err)
		return reply, err
	}
	return reply, nil
}
