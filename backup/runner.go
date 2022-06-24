package backup

import (
	"log"
	"sync"
)

const lineLength = 80

type Runner struct {
	Waiter     sync.WaitGroup
	configName string
	config     *configuration
	Stats      *stats
}

func NewRunner(configName string) *Runner {
	runner := &Runner{
		configName: configName,
	}

	config := getConfiguration(configName)

	runner.config = config
	runner.Waiter.Add(1)

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
