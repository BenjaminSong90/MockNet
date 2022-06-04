package utils

import (
	"log"
	"os"
	"os/signal"
)

func ListenBreak(breakFunc func()) {

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	breakFunc()
	log.Println("Receive Interrupt signal")
}
