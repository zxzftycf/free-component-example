package base

import (
	"errors"
	"math/rand"
	"net"
	"strings"
	"time"

	"../common"
	"github.com/gomodule/redigo/redis"
)

func init() {
	common.AllComponentMap["Find"] = &Find{}
}

type Find struct {
	Base
	RedisPool     redis.Pool
	serverIp      string
	registerTimer *time.Ticker
	grpcPort      string
}

func (self *Find) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)

	self.RedisPool = redis.Pool{
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
	conn, err := net.Dial("udp", "www.google.com.hk:80")
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	self.serverIp = strings.Split(conn.LocalAddr().String(), ":")[0]
	self.grpcPort = (*self.Config)["grpc_port"]
	common.LogInfo("FindComponent serverIp : ", self.serverIp, self.grpcPort)
	return
}

func (self *Find) RegisterComponent() {
	self.registerTimer = time.NewTicker(500 * time.Millisecond)

	go func(t *time.Ticker) {
		for {
			<-t.C
			conn := self.RedisPool.Get()
			defer conn.Close()
			curServerConfig := common.ServerConfig[common.ServerName]
			conn.Send("MULTI")
			for componentName, _ := range curServerConfig {
				componentInfo := componentName + "_" + self.serverIp + ":" + self.grpcPort
				conn.Send("SETEX", componentInfo, 1, "1")
			}
			_, err := conn.Do("EXEC")
			if err != nil {
				common.LogError("RegisterComponent op redis has err", err)
			}
			conn.Close()
			//activeCount := self.RedisPool.ActiveCount()
			//common.LogInfo("RegisterComponent", activeCount)
			/*for componentName, _ := range curServerConfig {
				componentInfo := componentName + "_" + self.serverIp + ":" + self.grpcPort
				_, err := conn.Do("SETEX", componentInfo, 1, "1")
				if err != nil {
					common.LogError("RegisterComponent op redis has err", componentInfo, err)
				}
			}*/
		}
	}(self.registerTimer)
}

func (self *Find) FindComponent(componentName string) (string, error) {
	conn := self.RedisPool.Get()
	defer conn.Close()
	rst, err := redis.Strings(conn.Do("KEYS", componentName+"*"))
	if err != nil {
		return "", err
	}
	if len(rst) <= 0 {
		return "", errors.New("Find FindComponent not find this component:" + componentName)
	}
	n := rand.Intn(len(rst))
	ip := strings.Split(rst[n], "_")[1]
	return ip, nil
}

func (self *Find) FindAllComponent(componentName string) (map[string]bool, error) {
	conn := self.RedisPool.Get()
	defer conn.Close()
	ips := make(map[string]bool)
	rst, err := redis.Strings(conn.Do("KEYS", componentName+"*"))
	if err != nil {
		return ips, err
	}
	if len(rst) <= 0 {
		return ips, nil
	}

	for _, ip := range rst {
		realIp := strings.Split(ip, "_")[1]
		ips[realIp] = true
	}

	return ips, nil
}
