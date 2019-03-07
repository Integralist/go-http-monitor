package alarms

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/integralist/go-http-monitor/internal/instrumentator"
)

// Alarm contains fields relevant to an exceeded monitoring threshold.
type Alarm struct {
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

	// sleepInterval := time.Minute * time.Duration(evaluation)
	iteration := 0
	evaluationSecs := 60 * evaluation

	for {
		// don't bother checking alarm threshold exceeded on program start
		if iteration == 0 {
			iteration = 1
			// time.Sleep(sleepInterval)
			time.Sleep(time.Second * time.Duration(30))
			continue
		}

		// lineCount, err := lineCounter(f)
		// if err != nil {
		// 	// TODO: use tags for raw error message
		// 	instr.Logger.Error(err)
		// }

		fileScanner := bufio.NewScanner(f)
		lineCount := 0
		for fileScanner.Scan() {
			lineCount++
		}

		// reset position back to zero
		f.Seek(0, 0)

		avg := float64(evaluationSecs) / float64(lineCount)

		if avg > float64(threshold) {
			alarmChannel <- Alarm{}
		}

		// time.Sleep(sleepInterval)
		time.Sleep(time.Second * time.Duration(30))
		continue
	}
}

// Alternative line counter that benefits from assembly optimized functions
// offered by the bytes package to search characters in a byte slice.
func lineCounter(r io.Reader) (int, error) {
	// my log lines are ~80 bytes in length
	buf := make([]byte, 80*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// Process an alarm sent to the specified channel
func Process(alarmChannel <-chan Alarm, instr *instrumentator.Instr) {
	instr.Logger.Debug("ALARM_PROCESSING")

	for a := range alarmChannel {
		fmt.Printf("alarm! %+v\n", a)
	}
}
