package internal

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WatchForSignal(signalVar *GoSafeVar[bool]) {
	sc := make(chan os.Signal, 1)

	//catch sigint, sigterm & os.int
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc

	// Sleep one minute to give the grader a little more time
	time.Sleep(1 * time.Minute)

	SetValue(signalVar, true)

}
