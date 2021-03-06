package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"sync"

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
	populate      *bool
	sections      []string
	statsInterval int
	threshold     int
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
	populate = flag.Bool("populate", false, "populate access log with simulated http requests")
	const (
		flagEvaluationValue    = 2
		flagEvaluationUsage    = "alarm monitoring evaluation period in minutes"
		flagLocationValue      = "./access.log"
		flagLocationUsage      = "location of access.log file to monitor"
		flagStatsIntervalValue = 10
		flagStatsIntervalUsage = "statistic output interval in seconds"
		flagThresholdValue     = 10
		flagThresholdUsage     = "alarm threshold for total number of requests on avg"
	)
	flag.IntVar(&evaluation, "evaluation", flagEvaluationValue, flagEvaluationUsage)
	flag.IntVar(&evaluation, "e", flagEvaluationValue, flagEvaluationUsage+" (shorthand)")
	flag.StringVar(&location, "location", flagLocationValue, flagLocationUsage)
	flag.StringVar(&location, "l", flagLocationValue, flagLocationUsage+" (shorthand)")
	flag.IntVar(&statsInterval, "stats", flagStatsIntervalValue, flagStatsIntervalUsage)
	flag.IntVar(&statsInterval, "s", flagStatsIntervalValue, flagStatsIntervalUsage+" (shorthand)")
	flag.IntVar(&threshold, "threshold", flagThresholdValue, flagThresholdUsage)
	flag.IntVar(&threshold, "t", flagThresholdValue, flagThresholdUsage+" (shorthand)")
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
	sections = []string{
		"foo",
		"bar",
		"baz",
		"qux",
		"qiz",
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

	// open a readonly file descriptor for the access log
	f, err := os.Open(location)
	if err != nil {
		instr.Logger.Fatal("ACCESS_OPEN_FAILED")
	}

	// open a readonly file descriptor for counting the access log. the rationale
	// for having another descriptor is because the Seek (and other File calls)
	// done within the processor package will mess up the line counting we need
	// to do within the alarms package.
	fc, err := os.Open(location)
	if err != nil {
		instr.Logger.Fatal("ACCESS_OPEN_FAILED")
	}

	// open a r/w file descriptor for the access log for dynamic population
	fileRW, err := os.Create(location)
	if err != nil {
		instr.Logger.Fatal("ACCESS_OPEN_FAILED")
	}

	// handle Ctrl-C from user to stop the program
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func(f io.ReadCloser, fileRW io.WriteCloser) {
		for sig := range sigs {
			// clean-up resources in case of failure
			f.Close()
			fileRW.Close()

			instr.Logger.Info("CTRL-C RECEIVED")
			instr.Logger.Info(sig)
			os.Exit(2)
		}
	}(f, fileRW)

	// channel creation for synchronizing data
	alarmChannel := make(chan alarms.Alarm)
	statChannel := make(chan stats.Stat)

	// initialize a thread safe map for passing into our stats processor
	statsTracking := new(sync.Map)

	// start various background goroutines, passing in their dependencies
	go processor.Process(f, statChannel, statsInterval, &instr)
	go alarms.Monitor(fc, alarmChannel, evaluation, threshold, &instr)
	go alarms.Process(alarmChannel, &instr)
	go stats.Process(statChannel, statsTracking, &instr)

	// keep program running until user stops it with <Ctrl-C>
	for {
		// if requested, we'll populate the given access log with simulated http requests
		if *populate {
			line := generator.RandomRequest(ips, usernames, pages, sections, generator.LastDate)
			generator.Generate(fileRW, line, &instr)
		}
	}
}
