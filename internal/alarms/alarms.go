package alarms

import (
	"github.com/integralist/go-http-monitor/internal/instrumentator"
)

// Alarm contains fields relevant to an exceeded monitoring threshold.
type Alarm struct {
}

// Process ...
func Process(alarmChannel chan Alarm, instr *instrumentator.Instr) {
	instr.Logger.Debug("ALARM_PROCESSING")
}
