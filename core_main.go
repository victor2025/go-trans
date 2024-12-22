package main

import (
	"C"
	"go-trans/handlers"
)

//export StartUp
func StartUp(port, basePath *C.char) {
	sHandler := handlers.NewReceiveHandler(C.GoString(port), C.GoString(basePath))
	go sHandler.Handle()
}

func main() {
	sHandler := handlers.NewReceiveHandler("20235", "")
	sHandler.Handle()
}
