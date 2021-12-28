package japanvisacovidbot

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	GlobalSignalHandler chan os.Signal
)

func LoadGlobalSignalHandler() {
	GlobalSignalHandler = make(chan os.Signal)
	signal.Notify(GlobalSignalHandler, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
}
