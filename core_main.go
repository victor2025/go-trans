package main

import (
	"C"
	"encoding/json"
	"go-trans/app"
	"go-trans/utils"
)

//export StartupWithConfig
func StartupWithConfig(configStr string) {
	var config map[string]any
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
