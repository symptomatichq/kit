package graceful

import (
	"os"
	"os/signal"
	"syscall"
)

type ShutdownFunc func(os.Signal)

func Handle(onExit ShutdownFunc) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		s := <-signals
		onExit(s)
	}()
}
