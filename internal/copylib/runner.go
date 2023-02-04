package copylib

import (
	"log"
	"sync"
)

type LogMode int8

const (
	LogSilent LogMode = iota
	LogSimple
	LogVerbose
)

var currentLogMode = LogSilent

const lineLength = 80

type Runner struct {
	Waiter     sync.WaitGroup
	configName string
	config     *configuration
	Stats      *stats
}

func NewRunner(configName string, logMode LogMode) *Runner {
	runner := &Runner{
		configName: configName,
	}

	config := getConfiguration(configName)

	runner.config = config
	runner.Waiter.Add(1)

	// set the current logging mode to what the user chose
	currentLogMode = logMode

	return runner
}

func (runner *Runner) Copy() {
	defer runner.handleFinish()

	fileCopier := &fileCopier{}
	fileCopier.run(runner.config)

	runner.Stats = &fileCopier.stats

	Print("file copy complete...")
}

func (runner *Runner) handleFinish() {
	recovery := recover()
	if recovery != nil {
		log.Printf("panic occurred:\n    %v", recovery)
	}

	runner.Waiter.Done()
}
