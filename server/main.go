package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"./base"
	"./common"
	"./logic"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	//读取服务编排配置
	file, err := os.Open("./layout.json")
	if err != nil {
		common.LogError("open layout.json has err", err)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	file.Close()
	err = json.Unmarshal(content, &common.ServerConfig)
	if err != nil {
		common.LogError(" json.Unmarshal has err", err)
		return
	}
	common.LogInfo("common.ServerConfig", common.ServerConfig)
	//读取服务器环境变量，加载组件
	serverName := os.Getenv("SERVER_NAME")
	if len(os.Args) > 1 {
		serverName = os.Args[1]
	}
	common.ServerName = serverName
	common.LogInfo("common.ServerName:", common.ServerName)
	common.ServerIndex = "0"
	serverIndex := os.Getenv("SERVER_INDEX")
	if len(os.Args) > 2 {
		serverIndex = os.Args[2]
	}
	if serverIndex != "" {
		common.ServerIndex = serverIndex
	}
	common.LogInfo("common.ServerIndex:", common.ServerIndex)
	curServerConfig := common.ServerConfig[serverName]
	if curServerConfig == nil {
		common.LogError("curServerConfig == nil")
		return
	}
	commonServerConfig := common.ServerConfig["common_config"]
	if commonServerConfig == nil {
		common.LogError("commonServerConfig == nil")
		return
	}
	mustServerConfig := common.ServerConfig["must"]
	if mustServerConfig == nil {
		common.LogError("mustServerConfig == nil")
		return
	}
	base.Init()
	logic.Init()
	//先加载基础组件
	for componentName, componentConfig := range mustServerConfig {
		if common.ComponentMap[componentName] != nil {
			continue
		}
		if common.AllComponentMap[componentName] == nil {
			common.LogError("init component err, componentName == nil", componentName)
			return
		}
		oneComponentConfig := common.OneComponentConfig{}
		//先用公共配置的值来填充
		commonComponentConfig := commonServerConfig[componentName]
		if commonComponentConfig != nil {
			for oneFieldName, oneFieldValue := range commonComponentConfig {
				oneComponentConfig[oneFieldName] = oneFieldValue
			}
		}
		//然后使用基础组件配置填充
		for oneFieldName, oneFieldValue := range componentConfig {
			oneComponentConfig[oneFieldName] = oneFieldValue
		}
		//最后使用各自进程的组件配置填充
		curComponentConfig := curServerConfig[componentName]
		if curComponentConfig != nil {
			for oneFieldName, oneFieldValue := range curComponentConfig {
				oneComponentConfig[oneFieldName] = oneFieldValue
			}
		}
		methodArgs := []reflect.Value{reflect.ValueOf(&oneComponentConfig), reflect.ValueOf(componentName)}
		reflect.ValueOf(common.AllComponentMap[componentName]).MethodByName("LoadComponent").Call(methodArgs)
		common.ComponentMap[componentName] = common.AllComponentMap[componentName]
	}
	//开始加载进程独有组件
	for componentName, componentConfig := range curServerConfig {
		if common.ComponentMap[componentName] != nil {
			continue
		}
		if common.AllComponentMap[componentName] == nil {
			common.LogError("init component err, componentName == nil", componentName)
			return
		}
		oneComponentConfig := common.OneComponentConfig{}
		//先用公共配置的值来填充
		commonComponentConfig := commonServerConfig[componentName]
		if commonComponentConfig != nil {
			for oneFieldName, oneFieldValue := range commonComponentConfig {
				oneComponentConfig[oneFieldName] = oneFieldValue
			}
		}
		//然后使用自己的配置填充
		for oneFieldName, oneFieldValue := range componentConfig {
			oneComponentConfig[oneFieldName] = oneFieldValue
		}
		methodArgs := []reflect.Value{reflect.ValueOf(&oneComponentConfig), reflect.ValueOf(componentName)}
		reflect.ValueOf(common.AllComponentMap[componentName]).MethodByName("LoadComponent").Call(methodArgs)
		common.ComponentMap[componentName] = common.AllComponentMap[componentName]
	}

	//开启服务发现与注册服务
	findComponentInterface := common.ComponentMap["Find"]
	if findComponentInterface != nil {
		findComponent, ok := findComponentInterface.(*base.Find)
		if !ok {
			common.LogError(" findComponentInterface not findComponent ")
			return
		}
		findComponent.RegisterComponent()
	}
	//开启分布式锁组件
	lockComponentInterface := common.ComponentMap["Lock"]
	if lockComponentInterface != nil {
		lockComponent, ok := lockComponentInterface.(*base.Lock)
		if !ok {
			common.LogError(" lockComponentInterface not lockComponent ")
			return
		}
		common.Locker = lockComponent
	}
	//开启消息路由组件
	routeComponentInterface := common.ComponentMap["Route"]
	if routeComponentInterface != nil {
		routeComponent, ok := routeComponentInterface.(*base.Route)
		if !ok {
			common.LogError(" routeComponentInterface not routeComponent ")
			return
		}
		common.Router = routeComponent
	}
	//开启推送组件
	pushComponentInterface := common.ComponentMap["Push"]
	if pushComponentInterface != nil {
		pushComponent, ok := pushComponentInterface.(*base.Push)
		if !ok {
			common.LogError(" pushComponentInterface not pushComponent ")
			return
		}
		common.Pusher = pushComponent
	}

	time.Sleep(time.Duration(2) * time.Second)
	//开启日志服务
	logComponentInterface := common.ComponentMap["Log"]
	if logComponentInterface != nil {
		logComponent, ok := logComponentInterface.(*base.Log)
		if !ok {
			common.LogError(" logComponentInterface not logComponent ")
			return
		}
		common.Logger = logComponent
	}

	common.LogInfo("server start ok", common.ServerName, common.ServerIndex)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGINT,
		syscall.SIGILL,
		syscall.SIGFPE,
		syscall.SIGSEGV,
		syscall.SIGTERM,
		syscall.SIGABRT)
	<-signalChan
	common.LogInfo("do some close operate")
	common.LogInfo("server end")
}
