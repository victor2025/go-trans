package test

import (
	"fmt"
	"go-trans/pkg/models/entity"
	"go-trans/pkg/orm"
	"go-trans/services"
	"go-trans/utils"
	"testing"
	"time"
)

/**
  @author: victor2022
  @since: 2025/1/12
*/

func TestDDL(t *testing.T) {
	db := orm.GetDb("../res/db/dev.db")
	db.Create(&entity.SendTaskInfo{
		BaseTaskInfo: &entity.BaseTaskInfo{
			BaseModel: &entity.BaseModel{
				GmtCreate: time.Now(),
				GmtModify: time.Now(),
			},
			FileName: "test.txt",
			FilePath: "test.txt",
			TaskId:   "002",
		},
		ReceiverId: "0",
	})
	sendTaskInfo := &entity.SendTaskInfo{}
	db.First(sendTaskInfo, "task_id = ?", "002")
	fmt.Println(sendTaskInfo)

	sendTaskInfo.Progress = 0.000231
	db.Save(sendTaskInfo)
	db.First(sendTaskInfo, "task_id = ?", "002")
	fmt.Println(sendTaskInfo)
}

func TestCreateDevice(t *testing.T) {
	db := orm.GetDb("../res/db/dev.db")
	deviceInfo := entity.GetNewDeviceInfo("localhost", "20235")
	deviceInfo.DeviceId = "001"
	deviceInfo.DeviceName = "LOCAL"
	db.Save(deviceInfo)
	sendTaskInfo := &entity.DeviceInfo{}
	db.First(sendTaskInfo, "device_id = ?", "001")
	fmt.Println(sendTaskInfo)

}

func TestCreateSendTaskInfo(t *testing.T) {
	db := services.GetServiceContext().DB
	err := services.GetServiceContext().TaskService.CreateSendTask("../bin", "002")
	utils.HandleError(err)

	result := &entity.SendTaskInfo{}
	db.First(result, 1)
	fmt.Println(result)

	result.Progress = 0.1
	services.GetServiceContext().TaskService.UpdateSendTask(result)
	result = &entity.SendTaskInfo{}
	db.First(result, 1)
	fmt.Println(result)
}
