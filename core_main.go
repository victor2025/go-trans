package main

import (
	"C"
	"go-trans/app"
)

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
