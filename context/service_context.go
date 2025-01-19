package context

import (
	goContext "context"
	"github.com/mustafaturan/bus/v3"
	"go-trans/pkg/models/dto"
	"go-trans/pkg/orm"
	"go-trans/services"
	"go-trans/utils"
	"gorm.io/gorm"
	"sync"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/

type ServiceContext struct {
	ServerService *services.TransmitService
	ConfigService *services.ConfigService
	TaskService   *services.SendTaskService
	DB            *gorm.DB
	BusService    *services.BusService
}

var serviceContext *ServiceContext
var once = sync.Once{}

// GetServiceContext 获取系统上下文
func GetServiceContext() *ServiceContext {
	once.Do(func() {
		// 配置服务
		configService := services.NewConfigService()
		// 持久层配置
		dbLocation := configService.GetOrDefault("orm.db.location", "./res/db/dev.db")
		DB := orm.GetDb(dbLocation)
		// 总线配置
		busTopic := configService.GetOrDefault("msg.bus.topic", "msg.bus")

		serviceContext = &ServiceContext{
			ServerService: services.NewServerService(),
			ConfigService: configService,
			TaskService:   services.NewTaskService(DB),
			BusService:    services.NewBusService(busTopic),
			DB:            DB,
		}

		// 启动发送任务处理器
		serviceContext.startSendTaskProcessor()
	})
	return serviceContext
}

func (s *ServiceContext) StartReceiveServer() {
	port := s.ConfigService.GetOrDefault("transmit.server.port", "20235")
	filePath := s.ConfigService.GetOrDefault("transmit.server.filepath", "./received")
	s.ServerService.StartReceiveServer(port, filePath)
}

func (s *ServiceContext) StopReceiveServer() {
	s.ServerService.StopReceiveServer()
}

func (s *ServiceContext) startSendTaskProcessor() {
	once.Do(func() {
		taskService := s.TaskService
		// 启动processor
		processor := services.NewSendTaskProcessor(func() ([]*dto.SendTaskDto, error) {
			return taskService.GetSendTaskDtos(1, 10, "")
		}, func(info *dto.SendTaskDto) {
			err := serviceContext.BusService.PostMsg(info)
			utils.HandleError(err)
		})

		go processor.Start()

		// 注册消息监听
		serviceContext.BusService.RegisterHandler("sendStatusMsgHandler", func(ctx goContext.Context, e bus.Event) {
			taskDto := e.Data.(*dto.SendTaskDto)
			err := taskService.UpdateSendTask(taskDto.Task)
			utils.HandleError(err)
		})
	})
}
