package alarms

import (
	"fmt"

	"github.com/integralist/go-http-monitor/internal/instrumentator"
)

// Alarm contains fields relevant to an exceeded monitoring threshold.
type Alarm struct {
}

// Process ...
func Process(alarmChannel <-chan Alarm, instr *instrumentator.Instr) {
	instr.Logger.Debug("ALARM_PROCESSING")

	for a := range alarmChannel {
		fmt.Printf("alarm! %+v\n", a)
	}
}
