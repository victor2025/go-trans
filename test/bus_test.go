package test

import (
	"context"
	"fmt"
	"github.com/mustafaturan/bus/v3"
	"go-trans/pkg/models/entity"
	context2 "go-trans/services"
	"go-trans/utils"
	"testing"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/
func TestBus(t *testing.T) {
	busService := context2.GetServiceContext().BusService
	busService.RegisterHandler("testHandler",
		func(ctx context.Context, event bus.Event) {
			fmt.Println(event.Data)
			fmt.Printf("data:%v\n", event.Data)
		})

	// send msg
	info, err := entity.GetNewSendTaskInfo("./bus_test.go", "111")
	utils.HandleError(err)
	err = busService.PostMsg(info)
	utils.HandleError(err)

	for {
	}
}
