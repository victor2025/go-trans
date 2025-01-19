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
	sendTaskInfoModel := &entity.SendTaskInfo{}
	if !db.Migrator().HasTable(sendTaskInfoModel) {
		err := db.AutoMigrate(sendTaskInfoModel)
		utils.HandleError(err, utils.ExitOnErr)
	}
	receiveTaskInfoModel := &entity.ReceiveTaskInfo{}
	if !db.Migrator().HasTable(receiveTaskInfoModel) {
		err := db.AutoMigrate(receiveTaskInfoModel)
		utils.HandleError(err, utils.ExitOnErr)
	}

	deviceInfoModel := &entity.DeviceInfo{}
	if !db.Migrator().HasTable(deviceInfoModel) {
		err := db.AutoMigrate(deviceInfoModel)
		utils.HandleError(err, utils.ExitOnErr)
	}
}
