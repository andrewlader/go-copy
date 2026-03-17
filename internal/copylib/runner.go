package copylib

import (
	"fmt"
	"sync"
)

type Runner struct {
	Waiter     sync.WaitGroup
	configName string
	config     *configuration
	Stats      *stats
}

func NewRunner(configName string) (*Runner, error) {
	config := getConfiguration(configName)
	if config == nil {
		return nil, fmt.Errorf("configuration with name '%s' not found", configName)
	}

	runner := &Runner{
		configName: configName,
	}

	runner.config = config
	runner.Waiter.Add(1)

	return runner, nil
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
		PrintError(fmt.Sprintf("panic occurred:\n    %v", recovery))
	}

	runner.Waiter.Done()
}
