package test

import (
	"context"
	"fmt"
	"github.com/mustafaturan/bus/v3"
	"go-trans/pkg/msg_bus"
	"go-trans/pkg/orm/models"
	"go-trans/utils"
	"testing"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/
func TestBus(t *testing.T) {
	msg_bus.RegisterHandler("testHandler",
		func(ctx context.Context, event bus.Event) {
			fmt.Println(event.Data)
			fmt.Printf("data:%v\n", event.Data)
		})

	// send msg
	info, err := models.GetNewSendTaskInfo("./bus_test.go", "111")
	utils.HandleError(err)
	err = msg_bus.PostMsg(info)
	utils.HandleError(err)

	for {
	}
}
