package stats

import (
	"fmt"
	"strings"
	"sync"

	"github.com/integralist/go-http-monitor/internal/formatter"
	"github.com/integralist/go-http-monitor/internal/instrumentator"
)

// Stat contains fields relevant to a statistical analysis.
type Stat struct {
	Logs []byte
}

// Analysis is a simplified version of sync.Map which will aid with testing.
type Analysis interface {
	Load(key interface{}) (value interface{}, ok bool)
	Range(f func(key, value interface{}) bool)
	Store(key, value interface{})
}

type requestedPages map[string]int
type requestData map[string]requestedPages

// Process data sent to the specified channel for statistical analysis.
func Process(
	statChannel <-chan Stat,
	statsTracking Analysis,
	instr *instrumentator.Instr) {

	instr.Logger.Debug("STAT_PROCESSING")

	// low-level primitive for ensuring data structure is thread-safe
	//
	// I decided to use this with a normal map data structure because sync.Map
	// have generic value types of interface{} that I'm not keen on.
	var mutex = &sync.Mutex{}

	for s := range statChannel {
		// data structures used for data analysis and tracking
		var counter int
		pages := make(requestedPages)
		data := make(requestData)

		for _, record := range strings.Split(string(s.Logs), "\n") {
			cells := strings.Split(record, " ")

			// the last record is just a line break
			if len(cells) < 10 {
				continue
			}

			// use variables to make identifying abstract data clearer
			ip := cells[0]
			user := cells[2]
			request := cells[6]

			// if ip is already tracked, let's extract the nested data structures so
			// we can increment the existing counter value for a given user/request.
			if d, ok := statsTracking.Load(ip); ok {
				// we need to type assert in order to interact with our nested map
				if rd, ok := d.(requestData); ok {
					data = rd

					// although the ip might already be tracked, the nested data
					// structure might not be initialized yet
					if data[user] == nil {
						data[user] = make(requestedPages)
					}
					pages = data[user]
					counter = pages[request]
				}
			}

			mutex.Lock()

			counter = counter + 1
			pages[request] = counter
			data[user] = pages

			mutex.Unlock()

			/*
				data structure we're ultimately ending up with:

					(statsTracking) map[ip]:
						(requestedData) map[username]:
						  (requestedPages) map[url]: counter
			*/
			statsTracking.Store(ip, data)
		}

		// we'll just dump out the users, and the number of times they requested
		// specific pages as a basic indicator of some stats collected.
		fmt.Println("--------------------------------------")
		fmt.Printf("Stats for last 10 seconds of requests:\n\n")
		formatter.Pretty(data)
		fmt.Println("--------------------------------------")
	}
}
