package services

import (
	"go-trans/pkg/models/consts"
	"log"
)

/**
  @author: victor2022
  @since: 2025/3/15
*/

func Info(v ...interface{}) {
	if !isLogEnabled() {
		return
	}
	log.Println(v...)
}

func InfoF(format string, v ...interface{}) {
	if !isLogEnabled() {
		return
	}
	log.Printf(format, v...)
}

func isLogEnabled() bool {
	flag := GetServiceContext().ConfigService.GetOrDefault(consts.LogEnable, consts.NO)
	return flag == consts.YES
}
