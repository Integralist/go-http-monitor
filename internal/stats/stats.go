package stats

import "github.com/integralist/go-http-monitor/internal/instrumentator"

// Process ...
func Process(instr *instrumentator.Instr) {
	instr.Logger.Debug("STATS_PROCESSING")
}
