package common

import (
	"fmt"

	pb "../grpc"
)

func LogInfo(a ...interface{}) {
	if Logger != nil {
		Logger.Info(a...)
	} else {
		fmt.Println(a...)
	}
}

func LogError(a ...interface{}) {
	if Logger != nil {
		Logger.Error(a...)
	} else {
		fmt.Println(a...)
	}
}

func LogDebug(a ...interface{}) {
	if Logger != nil {
		Logger.Debug(a...)
	} else {
		fmt.Println(a...)
	}
}

func GenShortId(idLen uint16, idNum uint16) []string {
	result := []string{}
	letterMap := make(map[string]bool)
	letter := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for _, v := range letter {
		letterMap[v] = true
	}
	for i := uint16(1); i <= idNum; i++ {
		oneId := ""
		letterNum := uint16(0)
		for oneLetter, _ := range letterMap {
			oneId = oneId + oneLetter
			letterNum++
			if letterNum >= idLen {
				result = append(result, oneId)
				break
			}
		}
	}
	return result
}

func GetGrpcErrorMessage(code pb.ErrorCode, message string) *pb.ErrorMessage {
	errMessage := &pb.ErrorMessage{}
	errMessage.Code = code
	errMessage.Message = message
	return errMessage
}
