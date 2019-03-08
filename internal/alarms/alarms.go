package alarms

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/integralist/go-http-monitor/internal/instrumentator"
)

// Alarm contains fields relevant to an exceeded monitoring threshold.
type Alarm struct {
	Message string
}

type readSeeker interface {
	io.Reader
	io.Seeker
}

// Monitor counts the number of lines in our access.log file and will trigger
// an alarm whenever the total number of requests reaches a predefined
// threshold across the specified evaluation period.
func Monitor(
	f readSeeker,
	alarmChannel chan<- Alarm,
	evaluation int,
	threshold int,
	instr *instrumentator.Instr) {

	var wasExceeded bool
	var lineTracker int

	sleepInterval := time.Minute * time.Duration(evaluation)
	iteration := 0
	evaluationSecs := 60 * evaluation

	for {
		// don't bother checking alarm threshold exceeded on program start
		if iteration == 0 {
			iteration = 1
			time.Sleep(sleepInterval)
			continue
		}

		currentLineCount := 0

		fileScanner := bufio.NewScanner(f)
		for fileScanner.Scan() {
			currentLineCount++
		}

		// reset position back to zero
		f.Seek(0, 0)

		numRequests := currentLineCount - lineTracker
		avg := float64(numRequests) / float64(evaluationSecs)
		lineTracker = currentLineCount

		fmt.Println("currentLineCount:", currentLineCount)
		fmt.Println("numRequests:", numRequests)
		fmt.Println("avg:", avg)

		now := time.Now().Format(time.UnixDate)

		if avg > float64(threshold) {
			wasExceeded = true
			msg := fmt.Sprintf("High traffic generated an alert - hits = %f, triggered at %s", avg, now)

			alarmChannel <- Alarm{
				Message: msg,
			}
		}

		if wasExceeded && avg < float64(threshold) {
			wasExceeded = false
			msg := fmt.Sprintf("High traffic alarm has now recovered at %s", now)

			alarmChannel <- Alarm{
				Message: msg,
			}
		}

		time.Sleep(sleepInterval)
		continue
	}
}

// Process an alarm sent to the specified channel
func Process(alarmChannel <-chan Alarm, instr *instrumentator.Instr) {
	instr.Logger.Debug("ALARM_PROCESSING")

	for a := range alarmChannel {
		fmt.Printf("alarm! %+v\n", a)
	}
}
