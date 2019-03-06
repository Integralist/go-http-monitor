package generator

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/integralist/go-http-monitor/internal/instrumentator"
)

// LastDate tracks the last date we generated for the access.log simulation
var LastDate time.Time

// Generate produces random access.log records
func Generate(f io.WriteCloser, line []byte, instr *instrumentator.Instr) {
	_, err := f.Write(line)
	if err != nil {
		instr.Logger.Error("ACCESS_WRITE_FAILED")
	}

	// note: would a time.Tick channel been more appropriate? 🤔
	time.Sleep(time.Millisecond * time.Duration(100))
}

// RandomRequest generates a fake access log record
func RandomRequest(ips, usernames, pages, sections []string, lastDate time.Time) []byte {
	ip := ips[rand.Intn(5)]
	username := usernames[rand.Intn(5)]
	pageIndex := rand.Intn(25)
	sectionIndex := rand.Intn(5)
	page := fmt.Sprintf("%s/%s", sections[sectionIndex], pages[pageIndex])
	contentLength := pageIndex + 1

	timeAdd := []time.Duration{
		time.Minute,
		time.Second,
	}

	var t time.Time

	// if no date set previously (we check against the type's zero value)
	if lastDate.String() == "0001-01-01 00:00:00 +0000 UTC" {
		t = time.Now().Local().Add(timeAdd[rand.Intn(2)] * time.Duration(1))
	} else {
		t = lastDate.Add(timeAdd[rand.Intn(2)] * time.Duration(1))
	}

	date := fmt.Sprintf("%02d/%s/%d:%02d:%02d:%02d +0000", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())

	// update 'public' variable, which will be passed back into this function
	LastDate = t

	tmpl := "%s - %s [%s] \"GET /%s HTTP/1.1\" 200 %d\n"
	req := fmt.Sprintf(tmpl, ip, username, date, page, contentLength)

	return []byte(req)
}
