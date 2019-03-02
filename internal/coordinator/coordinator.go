package coordinator

import (
	"github.com/fsnotify/fsnotify"
	"github.com/integralist/go-http-monitor/internal/instrumentator"
	"github.com/integralist/go-http-monitor/internal/thresholds"
)

// Coordinate is the mediator of the program.
func Coordinate(location string, a chan thresholds.Alarm, instr *instrumentator.Instr) {
	instr.Logger.Info("START COORDINATING")

	go access(location, instr)
	go alarms(a, instr)
	go stats(instr)
}

func access(location string, instr *instrumentator.Instr) {
	instr.Logger.Info("READ ACCESS LOG")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		instr.Logger.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					instr.Logger.Debug("FILE_MODIFIED")
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
				instr.Logger.Error("FILE_READ_ERROR")
			}
		}
	}()

	err = watcher.Add(location)
	if err != nil {
		instr.Logger.Fatal(err)
	}

	<-done
}

func alarms(a chan thresholds.Alarm, instr *instrumentator.Instr) {
	instr.Logger.Info("HANDLE ALARMS")
}

func stats(instr *instrumentator.Instr) {
	instr.Logger.Info("HANDLE STATS")
}
