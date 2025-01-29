package main

import (
	"C"
	"encoding/json"
	"go-trans/app"
	"go-trans/utils"
	"log"
)

//export StartupWithConfig
func StartupWithConfig(configStrC *C.char) {
	var config map[string]any
	configStr := C.GoString(configStrC)
	log.Printf("receive start up config:%v\n", configStr)
	err := json.Unmarshal([]byte(configStr), &config)
	utils.HandleError(err, utils.PanicOnError)
	app.RunWithConfig(config)
}

//export Startup
func Startup() {
	app.Run()
}

//export Shutdown
func Shutdown() {
	app.Shutdown()
}

func main() {
	app.Run()
}
