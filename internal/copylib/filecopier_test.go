package copylib

import (
	"io/fs"
	"testing"
)

func TestCopySuccess(t *testing.T) {
	var configName = "foo"
	var source = "f:\\games\\foobar\\saves"
	var destinations = []string{"g:\\game_backups\\foobar\\saves", "h:\\game_backups\\foobar\\saves", "i:\\game_backups\\foobar\\saves"}

	SetupTestFileSystemFunctions(destinations)
	defer func() { ReinitializeFileSystemFunctions() }()

	config := &configuration{
		name:         configName,
		source:       source,
		destinations: destinations,
		replace:      replaceSkipIfSame,
	}

	runner := &Runner{
		configName: config.name,
		config:     config,
	}

	runner.Waiter.Add(1)
	currentLogMode = LogVerbose

	createSimpleTestFiles()

	runner.Copy()
}

func TestCopyWithSkipSuccess(t *testing.T) {
	var configName = "foo"
	var source = "f:\\games\\foobar\\saves"
	var destinations = []string{"g:\\game_backups\\foobar\\saves", "h:\\game_backups\\foobar\\saves", "i:\\game_backups\\foobar\\saves"}

	SetupTestFileSystemFunctions(destinations)
	defer func() { ReinitializeFileSystemFunctions() }()

	Stat = StatDestFileExistsSuccess

	config := &configuration{
		name:         configName,
		source:       source,
		destinations: destinations,
		replace:      replaceSkipIfSame,
	}

	runner := &Runner{
		configName: config.name,
		config:     config,
	}

	runner.Waiter.Add(1)
	currentLogMode = LogVerbose

	createSimpleTestFiles()

	runner.Copy()
}

func TestCopyWithReplaceSuccess(t *testing.T) {
	var configName = "foo"
	var source = "f:\\games\\foobar\\saves"
	var destinations = []string{"g:\\game_backups\\foobar\\saves", "h:\\game_backups\\foobar\\saves", "i:\\game_backups\\foobar\\saves"}

	SetupTestFileSystemFunctions(destinations)
	defer func() { ReinitializeFileSystemFunctions() }()

	Stat = StatDestFileExistsSuccess

	config := &configuration{
		name:         configName,
		source:       source,
		destinations: destinations,
		replace:      replaceAlways,
	}

	runner := &Runner{
		configName: config.name,
		config:     config,
	}

	runner.Waiter.Add(1)
	currentLogMode = LogVerbose

	createSimpleTestFiles()

	runner.Copy()
}

func TestCopyOpenFailure(t *testing.T) {
	var configName = "foo"
	var source = "f:\\games\\foobar\\saves"
	var destinations = []string{"g:\\game_backups\\foobar\\saves", "h:\\game_backups\\foobar\\saves", "i:\\game_backups\\foobar\\saves"}

	SetupTestFileSystemFunctions(destinations)
	defer func() { ReinitializeFileSystemFunctions() }()

	Open = OpenFailure

	config := &configuration{
		name:         configName,
		source:       source,
		destinations: destinations,
		replace:      replaceSkipIfSame,
	}

	runner := &Runner{
		configName: config.name,
		config:     config,
	}

	runner.Waiter.Add(1)
	currentLogMode = LogVerbose

	createSimpleTestFiles()

	runner.Copy()
}

func TestCopyFailure(t *testing.T) {
	var configName = "foo"
	var source = "f:\\games\\foobar\\saves"
	var destinations = []string{"g:\\game_backups\\foobar\\saves", "h:\\game_backups\\foobar\\saves", "i:\\game_backups\\foobar\\saves"}

	SetupTestFileSystemFunctions(destinations)
	defer func() { ReinitializeFileSystemFunctions() }()

	Copy = CopyFailure

	config := &configuration{
		name:         configName,
		source:       source,
		destinations: destinations,
		replace:      replaceSkipIfSame,
	}

	runner := &Runner{
		configName: config.name,
		config:     config,
	}

	runner.Waiter.Add(1)
	currentLogMode = LogVerbose

	createSimpleTestFiles()

	runner.Copy()
}

func createSimpleTestFiles() {
	testFiles = make([]fs.DirEntry, 0, 3)

	dirEntry := createDirEntry("foobar001.txt", 8600)
	testFiles = append(testFiles, dirEntry)

	dirEntry = createDirEntry("foobar002.txt", 8640)
	testFiles = append(testFiles, dirEntry)

	dirEntry = createDirEntry("foobar003.txt", 86400)
	testFiles = append(testFiles, dirEntry)
}
