package test

import (
	"fmt"
	"go-trans/pkg/orm"
	"go-trans/pkg/orm/models"
	"testing"
	"time"
)

/**
  @author: victor2022
  @since: 2025/1/12
*/

func TestDDL(t *testing.T) {
	db := orm.GetDb("../res/db/dev.db")
	db.Create(&models.SendTaskInfo{
		BaseTaskInfo: models.BaseTaskInfo{
			BaseModel: models.BaseModel{
				GmtCreate: time.Now(),
				GmtModify: time.Now(),
			},
			FileName: "test.txt",
			FilePath: "test.txt",
			TaskId:   "002",
		},
		ReceiverId: "0",
	})
	sendTaskInfo := &models.SendTaskInfo{}
	db.First(sendTaskInfo, "task_id = ?", "002")
	fmt.Println(sendTaskInfo)

	sendTaskInfo.Progress = 0.000231
	db.Save(sendTaskInfo)
	db.First(sendTaskInfo, "task_id = ?", "002")
	fmt.Println(sendTaskInfo)
}

func TestCreateSendTaskInfo(t *testing.T) {
	db := orm.GetDb("../res/db/dev.db")
	info, _ := models.GetNewSendTaskInfo("../bin", "001")
	fmt.Println(info)
	db.Save(info)
	result := &models.SendTaskInfo{}
	db.First(result, "task_id = ?", info.TaskId)
	fmt.Println(result)
}
