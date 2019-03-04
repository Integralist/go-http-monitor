package stats

import (
	"fmt"

	"github.com/integralist/go-http-monitor/internal/alarms"
	"github.com/integralist/go-http-monitor/internal/instrumentator"
)

// Stat contains fields relevant to a statistical analysis.
type Stat struct {
}

// Process data sent to the specified channel for statistical analysis.
func Process(
	statChannel <-chan Stat,
	alarmChannel chan<- alarms.Alarm,
	instr *instrumentator.Instr) {

	instr.Logger.Debug("STAT_PROCESSING")

	for s := range statChannel {
		fmt.Printf("stat! %+v\n", s)

		alarmChannel <- alarms.Alarm{}
	}
}
