package generator

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/integralist/go-http-monitor/internal/instrumentator"
)

// Generate produces random access.log records
func Generate(f io.WriteCloser, line []byte, instr *instrumentator.Instr) {
	_, err := f.Write(line)
	if err != nil {
		instr.Logger.Error("ACCESS_WRITE_FAILED")
	}

	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
}

// RandomRequest generates a fake access log record
func RandomRequest(ips, usernames, pages []string) []byte {
	ip := ips[rand.Intn(4)]
	username := usernames[rand.Intn(4)]
	pageIndex := rand.Intn(25)
	page := pages[pageIndex]
	contentLength := pageIndex + 1

	t := time.Now().Local().Add(time.Hour * time.Duration(1))
	date := fmt.Sprintf("%02d/%s/%d:%02d:%02d:00 +0000", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())

	tmpl := "%s - %s [%s] \"GET /%s HTTP/1.1\" 200 %d\n"
	req := fmt.Sprintf(tmpl, ip, username, date, page, contentLength)

	return []byte(req)
}
