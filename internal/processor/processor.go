package processor

import (
	"io"
	"time"

	"github.com/integralist/go-http-monitor/internal/alarms"
	"github.com/integralist/go-http-monitor/internal/instrumentator"
	"github.com/integralist/go-http-monitor/internal/stats"
)

// Process reads the access.log at a set interval and then generates a
// stats task to be processed and displayed to the user.
func Process(
	f io.Writer,
	statChannel chan<- stats.Stat,
	alarmChannel chan<- alarms.Alarm,
	statInterval int,
	instr *instrumentator.Instr) {

	instr.Logger.Debug("ACCESS_LOG_PROCESSING")

	for {
		alarmChannel <- alarms.Alarm{}
		statChannel <- stats.Stat{}
		time.Sleep(time.Second * time.Duration(statInterval))
		continue
	}
}
