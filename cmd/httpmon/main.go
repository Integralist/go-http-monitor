package main

import (
	"flag"
	"io"
	"os"
	"os/signal"

	"github.com/integralist/go-http-monitor/internal/alarms"
	"github.com/integralist/go-http-monitor/internal/generator"
	"github.com/integralist/go-http-monitor/internal/instrumentator"
	"github.com/integralist/go-http-monitor/internal/processor"
	"github.com/integralist/go-http-monitor/internal/stats"
	"github.com/sirupsen/logrus"
)

// instr contains pre-configured instrumentation tools
var instr instrumentator.Instr

// resources that will be passed around various package functions
var (
	evaluation    int
	help          *bool
	ips           []string
	location      string
	pages         []string
	statsInterval int
	threshold     int
	unit          string
	usernames     []string
	version       string // set via -ldflags in Makefile
)

func init() {
	// instrumentation
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true) // TODO: benchmark for performance implications

	// flag configuration
	help = flag.Bool("help", false, "show available command flags")
	const (
		flagEvaluationValue    = 2
		flagEvaluationUsage    = "monitoring evaluation period in minutes"
		flagLocationValue      = "./access.log"
		flagLocationUsage      = "location of access.log file to monitor"
		flagStatsIntervalValue = 10
		flagStatsIntervalUsage = "statistic output interval in seconds"
		flagThresholdValue     = 10
		flagThresholdUsage     = "average alarm threshold time period"
		flagUnitValue          = "second"
		flagUnitUsage          = "unit of time of the alarm threshold"
	)
	flag.IntVar(&evaluation, "evaluation", flagEvaluationValue, flagEvaluationUsage)
	flag.IntVar(&evaluation, "e", flagEvaluationValue, flagEvaluationUsage+" (shorthand)")
	flag.StringVar(&location, "location", flagLocationValue, flagLocationUsage)
	flag.StringVar(&location, "l", flagLocationValue, flagLocationUsage+" (shorthand)")
	flag.IntVar(&statsInterval, "stats", flagStatsIntervalValue, flagStatsIntervalUsage)
	flag.IntVar(&statsInterval, "s", flagStatsIntervalValue, flagStatsIntervalUsage+" (shorthand)")
	flag.IntVar(&threshold, "threshold", flagThresholdValue, flagThresholdUsage)
	flag.IntVar(&threshold, "t", flagThresholdValue, flagThresholdUsage+" (shorthand)")
	flag.StringVar(&unit, "unit", flagUnitValue, flagUnitUsage)
	flag.StringVar(&unit, "u", flagUnitValue, flagUnitUsage+" (shorthand)")
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

	// fake data for the sake of simulating http requests
	ips = []string{
		"127.0.0.1",
		"127.0.0.2",
		"127.0.0.3",
		"127.0.0.4",
		"127.0.0.5",
	}
	usernames = []string{
		"Bob",
		"Jane",
		"Lisa",
		"Mark",
		"Simon",
	}
	pages = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
}

func main() {
	if *help == true {
		flag.PrintDefaults()
		return
	}

	// note: I like log messages to be a bit more structured so I typically opt
	// for a format such as 'VERB_STATE' and 'NOUN_STATE' (as this makes searching
	// for errors within a log aggregator easier).
	//
	// note: I also typically prefer the "no news is good news" approach: which is
	// where you only log errors or warnings (not info/debug), as that makes
	// debugging easier as you don't have to filter out pointless messages about
	// things you already expected to happen, and the logs can instead focus on
	// surfacing all the _unexpected_ things that happened.
	instr.Logger.Debug("STARTUP_SUCCESSFUL")

	// open a r/w file descriptor for the access log so we can seek next position
	f, err := os.Create(location)
	if err != nil {
		instr.Logger.Fatal("ACCESS_OPEN_FAILED")
	}

	// handle Ctrl-C from user to stop the program
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func(f io.WriteCloser) {
		for sig := range sigs {
			f.Close()
			instr.Logger.Info("CTRL-C RECEIVED")
			instr.Logger.Info(sig)
			os.Exit(2)
		}
	}(f)

	// channel creation for synchronizing data
	alarmChannel := make(chan alarms.Alarm)
	statChannel := make(chan stats.Stat)

	// start various background goroutines
	go processor.Process(location, statChannel, alarmChannel, statsInterval, &instr)
	go alarms.Process(alarmChannel, &instr)
	go stats.Process(statChannel, &instr)

	// keep program running until user stops it with <Ctrl-C>
	for {
		// in the mean time, if no active access log provided, we'll simulate http requests
		if location == "./access.log" {
			line := generator.RandomRequest(ips, usernames, pages, generator.LastDate)
			generator.Generate(f, line, &instr)
		}
	}
}
