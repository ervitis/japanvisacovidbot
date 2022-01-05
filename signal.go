package japanvisacovidbot

import (
	"os"
	"os/signal"
	"syscall"
)

type (
	GlobalSignalHandlerData struct {
		Signals chan os.Signal
	}
)

var (
	GlobalSignalHandler GlobalSignalHandlerData
)

func LoadGlobalSignalHandler() {
	GlobalSignalHandler.Signals = make(chan os.Signal, 1)
	signal.Notify(GlobalSignalHandler.Signals, syscall.SIGTERM, os.Interrupt)
}
