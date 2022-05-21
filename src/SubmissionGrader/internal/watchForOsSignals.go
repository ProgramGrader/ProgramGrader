package internal

import (
	"os"
	"os/signal"
	"syscall"
)

func WatchForSignal(signalVar *GoSafeVar[bool]) {
	sc := make(chan os.Signal, 1)
	//catch sigint, sigterm & os.int
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc

	SetValue(signalVar, true)

}
