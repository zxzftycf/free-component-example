package base

import (
	"log"
	"os"
	"path"
	"sync"
	"time"

	"../common"
)

const (
	LogBaseDir = "./log/"
)

func init() {
	common.AllComponentMap["Log"] = &Log{}
}

type Log struct {
	common.LogI
	Base
	loggerFile  *os.File
	lCreateFile sync.Mutex
	logPath     string
	OpenDebug   bool
	Logger      *log.Logger
}

func (self *Log) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)
	self.OpenDebug = false
	if (*self.Config)["open_debug"] == "true" {
		self.OpenDebug = true
	}

	spath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	sBasePath := path.Join(spath, LogBaseDir, common.ServerName) + "/"
	if err := os.MkdirAll(sBasePath, os.ModePerm); err != nil {
		panic(err)
	}
	self.logPath = sBasePath
	self.createLoggerFile()
	return
}

func (self *Log) Info(a ...interface{}) {
	self.Logger.Println(a...)
}

func (self *Log) Error(a ...interface{}) {
	b := append([]interface{}{" ERROR: "}, a...)
	self.Logger.Println(b...)
}

func (self *Log) Debug(a ...interface{}) {
	if !self.OpenDebug {
		return
	}
	self.Logger.Println(a...)
}

func (self *Log) closeFile(f *os.File) {
	go func() {
		time.Sleep(time.Second * 30)
		f.Close()
	}()
}

//记录到文件
func (self *Log) createLoggerFile() {
	self.lCreateFile.Lock()
	defer self.lCreateFile.Unlock()

	stime := time.Now().Format("20060102150405")
	sname := ""
	for {
		sname = self.logPath + stime + "-" + common.ServerIndex + ".log"
		if self.loggerFile != nil {
			self.closeFile(self.loggerFile)
		}
		f, err := os.OpenFile(sname, os.O_CREATE|os.O_EXCL|os.O_RDWR, os.ModePerm)
		if err != nil {
			continue
		}
		self.loggerFile = f
		loggerTemp := log.New(f, "", log.LstdFlags)
		self.Logger = loggerTemp
		//标准输出重定向
		os.Stdout = f
		os.Stderr = f

		break
	}

	//定时换
	timeNow := time.Now()
	timeNext := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour()+1, 0, 0, 0, time.Local)
	tm := time.NewTimer(time.Second * time.Duration(timeNext.Unix()-timeNow.Unix()))
	go func() {
		<-tm.C
		tm.Stop()
		self.createLoggerFile()
	}()
}
