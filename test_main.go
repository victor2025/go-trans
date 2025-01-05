package main

import (
	"fmt"
	"go-trans/http"
	"go-trans/services"
)

/**
  @author: victor2022
  @since: 2025/1/2
*/

func main() {

	config := services.NewConfigService().GetOrDefault("http.server.port", "")
	fmt.Printf("config:%v\n", config)
	//context.GetSystemContext().StartReceiveServer()
	//for i := 0; i < 2; i++ {
	//	time.Sleep(1 * time.Second)
	//	log.Printf("wait: %d", i)
	//}
	//context.GetSystemContext().StopReceiveServer()
	//for {
	//}
	http.StartHttpServer()
}
