package processor

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/integralist/go-http-monitor/internal/instrumentator"
	"github.com/integralist/go-http-monitor/internal/stats"
)

type statReaderSeeker interface {
	io.Reader
	io.Seeker
	Stat() (os.FileInfo, error)
}

// TODO:
// - set var with time.Now so we can track last ten seconds in for loop
// - in for loop, look at current time.Now and compare to var of time.Now
// - if time is greater or equal to ten seconds, then send stat message

// Process reads the access.log at a set interval and then generates a
// stats task to be processed and displayed to the user.
func Process(
	f statReaderSeeker,
	statChannel chan<- stats.Stat,
	statInterval int,
	instr *instrumentator.Instr) {

	instr.Logger.Debug("ACCESS_LOG_PROCESSING")

	var cursor int64
	iteration := 0
	sleepInterval := time.Second * time.Duration(statInterval)

	// move cursor to end of file (read for next read)
	f.Seek(cursor, io.SeekEnd)

	// get current size of file, so we can calculate the diff in size
	size := fileSize(f, instr)
	fmt.Println("initial size", size)

	for {
		// don't bother sending a stat message on program start
		if iteration == 0 {
			iteration = 1
			time.Sleep(sleepInterval)
			continue
		}

		// calculate buffer size to read new content into
		bufferSize := fileSize(f, instr) - size
		fmt.Println("buffer size", bufferSize)

		// read from last file position to end of file for new records
		buffer := make([]byte, bufferSize)
		bytesRead, err := f.Read(buffer)
		if err != nil {
			instr.Logger.Warn("FILE_STAT_FAILED")
		}
		size = size + int64(bytesRead)
		fmt.Println("new size", size)

		/*
			127.0.0.2 - Lisa [04/March/2019:09:45:22 +0000] "GET /w HTTP/1.1" 200 23
			127.0.0.2 - Mark [04/March/2019:09:46:22 +0000] "GET /a HTTP/1.1" 200 1
			127.0.0.2 - Bob [04/March/2019:09:46:23 +0000] "GET /t HTTP/1.1" 200 20
			127.0.0.3 - Simon [04/March/2019:09:47:23 +0000] "GET /d HTTP/1.1" 200 4
		*/

		// send relevant information
		statChannel <- stats.Stat{}

		time.Sleep(sleepInterval)
		continue
	}
}

func fileSize(f statReaderSeeker, instr *instrumentator.Instr) int64 {
	stat, err := f.Stat()
	if err != nil {
		instr.Logger.Fatal("FILE_STAT_FAILED")
	}
	return stat.Size()
}
