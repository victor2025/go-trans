package services

import (
	"context"
	"github.com/mustafaturan/bus/v3"
	"go-trans/pkg/msg_bus"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/
type BusService struct {
	bus   *bus.Bus
	topic string
}

func NewBusService(topic string) *BusService {
	return &BusService{
		bus:   msg_bus.InitBus(topic),
		topic: topic,
	}
}

// PostMsg 发送消息
func (s *BusService) PostMsg(content interface{}) error {
	ctx := context.WithValue(context.Background(), bus.CtxKeyTxID, nil)
	return s.bus.Emit(ctx, s.topic, content)
}

// RegisterHandler 注册handler
func (s *BusService) RegisterHandler(handlerName string, function func(ctx context.Context, e bus.Event)) {
	handler := bus.Handler{
		Handle:  function,
		Matcher: s.topic,
	}
	s.bus.RegisterHandler(handlerName, handler)
}
