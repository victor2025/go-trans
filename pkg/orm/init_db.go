package orm

import (
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

/**
  @author: victor2022
  @since: 2025/1/12
*/
func GetDb(location string) *gorm.DB {
	err := utils.CreateBaseDirForFile(location)
	utils.HandleError(err)
	db, err := gorm.Open(sqlite.Open(location), &gorm.Config{})
	utils.HandleError(err, utils.ExitOnErr)
	log.Printf("Create db success, location:%s\n", location)

	// 初始化表
	initTables(db)
	return db
}

func initTables(db *gorm.DB) {
	// 发送任务
	sendTaskInfoModel := &entity.SendTaskInfo{}
	err := db.AutoMigrate(sendTaskInfoModel)
	utils.HandleError(err, utils.ExitOnErr)
	// 接收任务
	receiveTaskInfoModel := &entity.ReceiveTaskInfo{}
	err = db.AutoMigrate(receiveTaskInfoModel)
	utils.HandleError(err, utils.ExitOnErr)
	// 设备信息
	deviceInfoModel := &entity.DeviceInfo{}
	err = db.AutoMigrate(deviceInfoModel)
	utils.HandleError(err, utils.ExitOnErr)
}
