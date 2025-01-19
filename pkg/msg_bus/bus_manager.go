package msg_bus

import (
	"context"
	"github.com/mustafaturan/bus/v3"
	"github.com/mustafaturan/monoton/v2"
	"github.com/mustafaturan/monoton/v2/sequencer"
	"sync"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/
var msgBus *bus.Bus
var once sync.Once

const (
	msgBusTopic = "msg.bus"
)

func GetBus() *bus.Bus {
	once.Do(func() {
		initBus()
	})
	return msgBus
}

func initBus() {
	node := uint64(1)
	initialTime := uint64(1577865600000) // set 2020-01-01 PST as initial time
	m, err := monoton.New(sequencer.NewMillisecond(), node, initialTime)
	if err != nil {
		panic(err)
	}

	// init an id generator
	var idGenerator bus.Next = m.Next

	// create a new msg_bus instance
	b, err := bus.NewBus(idGenerator)
	if err != nil {
		panic(err)
	}

	// maybe register topics in here
	b.RegisterTopics(msgBusTopic, "order.fulfilled")

	msgBus = b
}

// RegisterHandler 注册handler
func RegisterHandler(handlerName string, function func(ctx context.Context, e bus.Event)) {
	handler := bus.Handler{
		Handle:  function,
		Matcher: msgBusTopic,
	}
	GetBus().RegisterHandler(handlerName, handler)
}

// PostMsg 发送消息
func PostMsg(content interface{}) error {
	ctx := context.WithValue(context.Background(), bus.CtxKeyTxID, nil)
	return GetBus().Emit(ctx, msgBusTopic, content)
}
