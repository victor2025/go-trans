package services

import (
	"go-trans/pkg/models/consts"
	"log"
)

/**
  @author: victor2022
  @since: 2025/3/15
*/

func Println(v ...interface{}) {
	if !isLogEnabled() {
		return
	}
	log.Println(v...)
}

func Printf(format string, v ...interface{}) {
	if !isLogEnabled() {
		return
	}
}

func isLogEnabled() bool {
	flag := GetServiceContext().ConfigService.GetOrDefault(consts.LogEnable, consts.NO)
	return flag == consts.YES
}
