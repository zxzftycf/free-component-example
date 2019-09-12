package base

import (
	"database/sql"
	"time"

	"../common"
	pb "../grpc"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

func init() {
	common.AllComponentMap["Mysql"] = &Mysql{}
}

type Mysql struct {
	Base
	Db *sql.DB
}

func (self *Mysql) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.Base.LoadComponent(config, componentName)
	dbTemp, err := sql.Open("mysql", (*self.Config)["connect_string"])
	if err != nil {
		panic(err)
	}
	self.Db = dbTemp
	err = self.checkTable()
	if err != nil {
		panic(err)
	}
	common.LogInfo("Mysql LoadComponent ok")
	return
}

func (self *Mysql) checkTable() error {
	_, err := self.Db.Exec(common.Mysql_Check_PlayerInfo_Table)
	if err != nil {
		common.LogError("Mysql checkTable Mysql_Check_PlayerInfo_Table has err", err)
		return err
	}
	_, err = self.Db.Exec(common.Mysql_Check_AccountInfo_Table)
	if err != nil {
		common.LogError("Mysql checkTable Mysql_Check_AccountInfo_Table has err", err)
		return err
	}
	return nil
}

func (self *Mysql) NewAccount(request *pb.MysqlAccountInfo, extroInfo *pb.MessageExtroInfo) (*pb.MysqlAccountInfo, error) {
	account := request.GetAccount()
	password := request.GetPassword()
	reply := &pb.MysqlAccountInfo{}
	var checkUuid string
	err := self.Db.QueryRow("select uuid from "+common.Mysql_AccountInfo_Table+" where account = ?", account).Scan(&checkUuid)
	if err == nil {
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_AccountExist, "")
		common.LogInfo("Mysql NewAccount account exist", account)
		return reply, nil
	}
	if err != nil && err != sql.ErrNoRows {
		common.LogError("Mysql NewAccount has err", err)
		return reply, err
	}
	userUuid, err := uuid.NewV4()
	if err != nil {
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_GenUuidFailed, "")
		common.LogInfo("Mysql NewAccount gen uuid fail", account)
		return reply, nil
	}
	//产生6位，10租短id
	nowTime := time.Now().Unix()
	shortIdArr := common.GenShortId(6, 10)
	for _, shortId := range shortIdArr {
		err := self.Db.QueryRow("select uuid from "+common.Mysql_AccountInfo_Table+" where short_id = ?", shortId).Scan(&checkUuid)
		if err == nil {
			continue
		}
		if err != nil && err != sql.ErrNoRows {
			common.LogError("Mysql NewAccount check short_id has err", err)
			return reply, err
		}
		_, err = self.Db.Exec("insert into "+common.Mysql_AccountInfo_Table+" (uuid,short_id,account,password,update_time) values(?,?,?,PASSWORD(?),?)", userUuid.String(), shortId, account, password, nowTime)
		if err != nil {
			common.LogError("Mysql NewAccount insert has err", err)
			return reply, err
		}
		reply.Uuid = userUuid.String()
		reply.ShortId = shortId
		reply.Account = account
		return reply, nil
	}

	reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_GenShortIdFailed, "")
	common.LogInfo("Mysql NewAccount gen short id fail", account)
	return reply, nil
}

func (self *Mysql) VerifyAccount(request *pb.MysqlAccountInfo, extroInfo *pb.MessageExtroInfo) (*pb.MysqlAccountInfo, error) {
	account := request.GetAccount()
	password := request.GetPassword()
	reply := &pb.MysqlAccountInfo{}
	err := self.Db.QueryRow("select uuid,short_id,update_time from "+common.Mysql_AccountInfo_Table+" where account = ? and password = PASSWORD(?)", account, password).Scan(&reply.Uuid, &reply.ShortId, &reply.UpdateTime)
	if err == sql.ErrNoRows {
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_DataNotFound, "")
		common.LogInfo("Mysql VerifyAccount account not found", account, password)
		return reply, nil
	}
	if err != nil {
		common.LogError("Mysql VerifyAccount has err", err)
		return reply, err
	}
	return reply, nil
}

func (self *Mysql) QueryKV(request *pb.MysqlKVMessage, extroInfo *pb.MessageExtroInfo) (*pb.MysqlKVMessage, error) {
	uuid := request.GetUuid()
	shortId := request.GetShortId()
	tableName := request.GetTableName()
	reply := &pb.MysqlKVMessage{}
	var err error
	if uuid != "" {
		err = self.Db.QueryRow("select info from "+tableName+" where uuid = ?", uuid).Scan(&reply.Info)
	} else if shortId != "" {
		err = self.Db.QueryRow("select info from "+tableName+" where short_id = ?", shortId).Scan(&reply.Info)
	}
	if err == sql.ErrNoRows {
		reply.ErrMessage = common.GetGrpcErrorMessage(pb.ErrorCode_DataNotFound, "")
		common.LogInfo("Mysql QueryKV data not found", uuid, shortId, tableName)
		return reply, nil
	}
	if err != nil {
		common.LogError("Mysql QueryKV has err", err)
		return reply, err
	}
	return reply, nil
}

func (self *Mysql) InsertKV(request *pb.MysqlKVMessage, extroInfo *pb.MessageExtroInfo) (*pb.MysqlKVMessage, error) {
	uuid := request.GetUuid()
	shortId := request.GetShortId()
	tableName := request.GetTableName()
	info := request.GetInfo()
	nowTime := time.Now().Unix()
	reply := &pb.MysqlKVMessage{}
	_, err := self.Db.Exec("insert into "+tableName+" (uuid,short_id,info,update_time) values(?,?,?,?)", uuid, shortId, info, nowTime)
	if err != nil {
		common.LogError("Mysql NewAccount insert has err", err)
		return reply, err
	}
	return reply, nil
}
