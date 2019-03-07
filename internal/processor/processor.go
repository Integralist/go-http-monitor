package processor

import (
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

	for {
		// don't bother sending a stat message on program start
		if iteration == 0 {
			iteration = 1
			time.Sleep(sleepInterval)
			continue
		}

		// calculate buffer size to read new content into
		bufferSize := fileSize(f, instr) - size

		// read from last file position to end of file for new records
		buffer := make([]byte, bufferSize)
		bytesRead, err := f.Read(buffer)
		if err != nil {
			// TODO: utilize tags for the raw error message and use FILE_STAT_FAILED
			// for the main message
			instr.Logger.Error(err)
		}
		size = size + int64(bytesRead)

		// send relevant log information
		//
		// TODO: the design is a bit odd here, we should construct a stat object
		// from within the stats package, so this should be sending just the logs
		// to the stats channel for processing instead.
		statChannel <- stats.Stat{
			Logs: buffer,
		}

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
