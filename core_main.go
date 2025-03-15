package main

import (
	"C"
	"encoding/json"
	"fmt"
	"go-trans/app"
	"go-trans/utils"
	"log"
)

//export StartupWithConfig
func StartupWithConfig(configStrC *C.char) *C.char {
	var config map[string]any
	configStr := C.GoString(configStrC)
	log.Printf("receive start up config:%v\n", configStr)
	err := json.Unmarshal([]byte(configStr), &config)
	utils.HandleError(err, utils.PanicOnError)
	port := app.RunWithConfig(config)
	return C.CString(fmt.Sprintf("%d", port))
}

//export Startup
func Startup() int {
	return app.Run()
}

//export Shutdown
func Shutdown() {
	app.Shutdown()
}

func main() {
	app.Run()
}
