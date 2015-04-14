package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/gotgo/lg"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TODO:
// LogConfig - load from json file, use config, reload config
// Log to Files - Crash, Alert, Main
// Global Values - Environment, AppName
// Graceful Shutdown
// Panic - log crash

var (
	configArg = flag.String("config", "./config.json", "path to config file")
	envArg    = flag.String("env", "dev", "operating environment (dev, stage, qa, prod)")
)

var (
	logMain *lg.LogrusReceiver
	logErr  *lg.LogrusReceiver
	logFile *lumberjack.Logger
)

func init() {
	logMain = &lg.LogrusReceiver{}
	logMain.Current().Formatter = new(logrus.TextFormatter)
	lg.AddReceiver(logMain)

	logErr = &lg.LogrusReceiver{}
	logErr.Current().Level = logrus.WarnLevel
	logErr.Current().Out = os.Stderr
	lg.AddReceiver(logErr)
}

func main() {
	flag.Parse()
	loadConfig() //will panic if we can load the config

	lg.Inform("Program Started", lg.KV{
		"app":    AppName,
		"rev":    CommitHash,
		"logger": "logrus"})

	reloadOnSIGHUP()

	lg.Inform("Program Ended")
}

func reloadOnSIGHUP() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			if logFile != nil {
				logFile.Rotate()
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
	if environment() == "prod" {
		logMain.Current().Formatter = new(logrus.JSONFormatter)
		logFile = &lumberjack.Logger{
			Filename:   "/var/log/myapp/foo.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		}
		logMain.Current().Out = logFile
	}
}

func environment() string {
	return *envArg
}
