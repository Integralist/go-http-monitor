package coordinator

import "github.com/integralist/go-http-monitor/internal/instrumentator"
import "github.com/integralist/go-http-monitor/internal/thresholds"

// Coordinate is the mediator of the program.
func Coordinate(a chan thresholds.Alarm, instr *instrumentator.Instr) {
	instr.Logger.Info("START COORDINATING")

	go access(instr)
	go alarms(a, instr)
	go stats(instr)
}

func access(instr *instrumentator.Instr) {
	instr.Logger.Info("READ ACCESS LOG")
}

func alarms(a chan thresholds.Alarm, instr *instrumentator.Instr) {
	instr.Logger.Info("HANDLE ALARMS")
}

func stats(instr *instrumentator.Instr) {
	instr.Logger.Info("HANDLE STATS")
}