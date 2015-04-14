package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gotgo/lg"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TODO:
// Panic - log crash
// Log to Files - Panic (stderr), Alert, Main
// LogConfig - load from json file, use config, reload config
// Global Values - Environment, AppName
// Graceful Shutdown

var (
	configArg   = flag.String("config", "./config.json", "path to config file")
	ctxArg      = flag.String("ctx", "dev", "operating context (local, dev, stage, qa, prod)")
	forceTTYArg = flag.Bool("tty", false, "output to stdout even in prod")
)

var (
	logMain  *lg.LogrusReceiver
	logAlert *lg.LogrusReceiver
	logFile  *lumberjack.Logger
	onClose  []func()
)

// regarding panics on seperate go routines -
// they will cause a crash and therefore can't be formatted into json
// (there's no recovery in the main thread of a panic in a go routine) and will get logged to stderr by go
func init() {
	logMain = &lg.LogrusReceiver{}
	logMain.Current().Formatter = new(logrus.TextFormatter)
	lg.AddReceiver(logMain)
	onClose = make([]func(), 0)
}

func main() {
	defer closeable()
	flag.Parse()

	loadConfig() //will panic if we can load the config

	lg.Inform("Program Started", lg.KV{
		"app":    AppName,
		"rev":    CommitHash,
		"logger": "logrus"})

	lg.Warn("trying out a warning")

	reloadOnSIGHUP()

	run(func() {
		time.Sleep(time.Second * 32)
	}, func() {
		time.Sleep(time.Second * 2)
	})

	lg.Inform("Program Ended")
}

func run(start, stop func()) {
	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGINT)

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer func() {
			wg.Done()
		}()
		lg.Inform("Running!")
		start()
	}()

	select {
	case <-term:
		lg.Inform("Got shutdown signal")
	}
	lg.Inform("Stopping... ")

	const tenSeconds = time.Second * 10

	startAt := time.Now()
	ticker := time.NewTicker(tenSeconds)
	go func() {
		for t := range ticker.C {
			seconds := int(t.Sub(startAt) / time.Second)
			lg.Warn("Still Waiting for Shutdown", lg.KV{"seconds": seconds})
		}
	}()

	stop()
	lg.Inform("Waiting on server to stop...")
	wg.Wait()
	ticker.Stop()
	lg.Inform("Done!")
}

func closeable() {
	for _, c := range onClose {
		if c != nil {
			c()
		}
	}
}

func reloadOnSIGHUP() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			if logFile != nil {
				logFile.Rotate() //?
			}
			loadConfig()
		}
	}()
}

func loadConfig() {
	//load config
	//config := &lg.LogConfig{}
}

func setupLogging() {
	ctx := *ctxArg
	if ctx != "dev" && ctx != "local" {
		logMain.Current().Formatter = new(logrus.JSONFormatter)
		logFile = &lumberjack.Logger{
			Filename:   "/var/log/myapp/main.json",
			MaxSize:    500, // megabytes
			MaxBackups: 4,
			MaxAge:     28, //days
		}
		logMain.Current().Out = logFile
		onClose = append(onClose, func() { logFile.Close() })

		// Put formattable Alerts Warn, Error and recovered Panic in the same file
		logAlert = &lg.LogrusReceiver{}
		logAlert.Current().Level = logrus.WarnLevel
		logAlert.Current().Formatter = new(logrus.JSONFormatter)
		//we discard on setup, in a prod server config, this is sent to a seperate file
		alertFile := &lumberjack.Logger{
			Filename:   "/var/log/myapp/alert.json",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		}
		logAlert.Current().Out = alertFile
		lg.AddReceiver(logAlert)
		onClose = append(onClose, func() { alertFile.Close() })

		//only applies in this context
		if *forceTTYArg {
			tty := &lg.LogrusReceiver{}
			tty.Current().Formatter = new(logrus.TextFormatter)
			lg.AddReceiver(logMain)
		}
	}
}
