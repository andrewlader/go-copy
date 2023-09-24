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

func TestCopyWithSubDirectoriesSuccess(t *testing.T) {
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

	createDirectoriesAndTestFiles()

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

	dirEntry := createDirEntry("foobar001.txt", 8600, false)
	testFiles = append(testFiles, dirEntry)

	dirEntry = createDirEntry("foobar002.txt", 8640, false)
	testFiles = append(testFiles, dirEntry)

	dirEntry = createDirEntry("foobar003.txt", 86400, false)
	testFiles = append(testFiles, dirEntry)
}

func createDirectoriesAndTestFiles() {
	testFiles = make([]fs.DirEntry, 0, 3)

	dirEntry := createDirEntry("foobar001.txt", 8600, false)
	testFiles = append(testFiles, dirEntry)

	dirEntry = createDirEntry("foobar002.txt", 8640, false)
	testFiles = append(testFiles, dirEntry)

	dirEntry = createDirEntry("foobar003.txt", 86400, false)
	testFiles = append(testFiles, dirEntry)

	dirEntrySubDir := createDirEntry("subdir001", 0, true)
	testFiles = append(testFiles, dirEntrySubDir)

	dirEntry = createDirEntry("foobar004.txt", 8600, false)
	dirEntrySubDir.addChildDirEntry(dirEntry)

	dirEntry = createDirEntry("foobar005.txt", 8640, false)
	dirEntrySubDir.addChildDirEntry(dirEntry)

	dirEntry = createDirEntry("foobar006.txt", 86400, false)
	dirEntrySubDir.addChildDirEntry(dirEntry)
}
