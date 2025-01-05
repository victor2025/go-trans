package main

import "go-trans/http"

/**
  @author: victor2022
  @since: 2025/1/2
*/

func main() {
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
