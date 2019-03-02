package main

import (
	"flag"
	"os"

	"github.com/integralist/go-http-monitor/internal/instrumentator"
	"github.com/sirupsen/logrus"
)

// instr contains pre-configured instrumentation tools
var instr instrumentator.Instr

var (
	location string
	version  string // set via -ldflags in Makefile
)

func init() {
	// instrumentation
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true) // TODO: benchmark for performance implications

	// flag configuration
	const (
		flagLocationValue = "./"
		flagLocationUsage = "location of access.log file to monitor"
	)
	flag.StringVar(&location, "location", flagLocationValue, flagLocationUsage)
	flag.StringVar(&location, "l", flagLocationValue, flagLocationUsage+" (shorthand)")
	flag.Parse()

	// instrumentation configuration
	//
	// we would in a real-world application configure this with additional fields
	// such as `Metric` (for handling the recording of metrics using a service
	// such as Datadog, just as an example)
	//
	// note: I prefer to configure instrumentation within the init function of
	// the main package, but because I'm then passing this struct instance around
	// to other functions in other packages, it means I need to use an exported
	// reference from a mediator package (i.e. the instrumentator package)
	instr = instrumentator.Instr{
		Logger: logrus.WithFields(logrus.Fields{
			"version":  version,
			"location": location,
		}),
	}
}

func main() {
	// note: I like log messages to be a bit more structured so I typically opt
	// for a format such as 'VERB_STATE' and 'NOUN_STATE' (as this makes searching
	// for errors within a log aggregator easier).
	//
	// note: I also typically prefer the "no news is good news" approach: which is
	// where you only log errors or warnings (not info/debug), as that makes
	// debugging easier as you don't have to filter out pointless messages about
	// things you already expected to happen, and the logs can instead focus on
	// surfacing all the _unexpected_ things that happened.
	instr.Logger.Info("STARTUP_SUCCESSFUL")
}
