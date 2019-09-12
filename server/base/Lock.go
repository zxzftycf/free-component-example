package base

import (
	"time"

	"../common"
	pb "../grpc"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

func init() {
	common.AllComponentMap["Lock"] = &Lock{}
}

type Lock struct {
	common.RedisLockI
	Base
	Redsync *redsync.Redsync
}

func (self *Lock) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)

	pools := []redsync.Pool{
		createPool((*self.Config)["redis_host"]),
	}

	self.Redsync = redsync.New(pools)
	return
}

func (self *Lock) Lock(name string) error {
	mutex := self.Redsync.NewMutex(name,
		redsync.SetExpiry(10*time.Second),
		redsync.SetTries(100),
		redsync.SetRetryDelay(100*time.Millisecond),
	)
	err := mutex.Lock()
	if err != nil {
		return err
	}
	return nil
}

func (self *Lock) Unlock(name string) {
	mutex := self.Redsync.NewMutex(name,
		redsync.SetExpiry(10*time.Second),
		redsync.SetTries(100),
		redsync.SetRetryDelay(100*time.Millisecond),
	)
	mutex.Unlock()
}

func (self *Lock) MessageLock(name string, extroInfo *pb.MessageExtroInfo, componentName string) error {
	if extroInfo.GetLocks() == nil {
		extroInfo.Locks = []*pb.MessageLock{}
	}
	for _, messageLock := range extroInfo.Locks {
		if messageLock.GetLockName() == name {
			newMessageLock := &pb.MessageLock{}
			newMessageLock.ComponentName = componentName
			newMessageLock.LockName = name
			newMessageLock.IsRealLock = false
			extroInfo.Locks = append(extroInfo.Locks, newMessageLock)
			common.LogDebug("Lock MessageLock old lock ok", extroInfo.Locks, componentName, name)
			return nil
		}
	}
	mutex := self.Redsync.NewMutex(name,
		redsync.SetExpiry(10*time.Second),
		redsync.SetTries(100),
		redsync.SetRetryDelay(100*time.Millisecond),
	)
	err := mutex.Lock()
	if err != nil {
		return err
	}
	newMessageLock := &pb.MessageLock{}
	newMessageLock.ComponentName = componentName
	newMessageLock.LockName = name
	newMessageLock.IsRealLock = true
	extroInfo.Locks = append(extroInfo.Locks, newMessageLock)
	common.LogDebug("Lock MessageLock new lock ok", extroInfo.Locks, componentName, name)
	return nil
}

func (self *Lock) MessageUnlock(name string, extroInfo *pb.MessageExtroInfo, componentName string) {
	if extroInfo.GetLocks() == nil || len(extroInfo.GetLocks()) <= 0 {
		common.LogError("Lock MessageUnlock has err extroInfo.GetLocks() == nil", componentName, name)
		return
	}
	allLocks := extroInfo.GetLocks()
	lastLock := allLocks[len(allLocks)-1]
	if lastLock.GetComponentName() != componentName || lastLock.GetLockName() != name {
		common.LogError("Lock MessageUnlock has err unlock component mismatch with lock component", extroInfo.Locks, componentName, name)
		return
	}
	if lastLock.GetIsRealLock() == false {
		extroInfo.Locks = append(allLocks[:len(allLocks)-1])
		common.LogDebug("Lock MessageUnlock not real lock ok", extroInfo.Locks, componentName, name)
		return
	}
	mutex := self.Redsync.NewMutex(name,
		redsync.SetExpiry(10*time.Second),
		redsync.SetTries(100),
		redsync.SetRetryDelay(100*time.Millisecond),
	)
	mutex.Unlock()
	extroInfo.Locks = append(allLocks[:len(allLocks)-1])
	common.LogDebug("Lock MessageUnlock real lock ok", extroInfo.Locks, componentName, name)
	return
}

func createPool(url string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 180 * time.Second, // Default is 300 seconds for redis server
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", url)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
