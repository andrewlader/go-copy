package copylib

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type copyContext struct {
	filename        string
	sourcePath      string
	destinationPath string
	subFolderPath   string
}

type stats struct {
	NumberOfSourceFiles  int
	NumberOfDestinations int
	TotalFilesSkipped    int
	TotalFilesCopied     int
	BytesCopied          int64
	TimeToCopy           time.Duration
}

type fileCopier struct {
	config *configuration
	stats  stats
}

func (fileCopier *fileCopier) run(config *configuration) {
	defer fileCopier.handleFinish()

	fileCopier.config = config
	fileCopier.stats.NumberOfDestinations = len(config.destinations)

	startTime := time.Now()
	fileCopier.walkPath("")
	fileCopier.stats.TimeToCopy = time.Now().Sub(startTime)
}

func (fileCopier *fileCopier) walkPath(pathToWalk string) {
	var file os.FileInfo

	context := &copyContext{
		sourcePath:    fileCopier.config.source,
		subFolderPath: pathToWalk,
	}

	currentPath := path.Join(fileCopier.config.source, context.subFolderPath)

	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		PrintWarning(fmt.Sprintf("skipping path %s: %v", currentPath, err))
	} else {
		for _, file = range files {
			if file.Mode().IsRegular() {
				// copy the file
				context.filename = file.Name()
				fileCopier.copyFileToDestinations(context)
			}
			if file.IsDir() {
				// walk the sub-folder path
				fileCopier.walkPath(path.Join(context.subFolderPath, file.Name()))
			}
		}
	}
}

func (fileCopier *fileCopier) copyFileToDestinations(context *copyContext) {
	var err error
	var count = 0
	var ok bool

	for _, destPath := range fileCopier.config.destinations {
		context.destinationPath = destPath

		// make sure all the sub folders exist for this destination path
		destinationPath := path.Join(context.destinationPath, context.subFolderPath)
		err = os.MkdirAll(destinationPath, os.ModeDir)
		if err != nil {
			// failed to create the sub-folder(s), so skip this path and continue
			continue
		}

		ok, err = fileCopier.copyFile(context, destinationPath)
		if err != nil {
			PrintError(fmt.Sprintf("error copying file %s: %s", context.filename, err))
		} else if ok {
			count++
		}
	}

	if count == 0 {
		Print(fmt.Sprintf("file \"%s\" was skipped", context.filename))
	} else if count == len(fileCopier.config.destinations) {
		Print(fmt.Sprintf("copied file \"%s\"", context.filename))
	} else {
		Print(fmt.Sprintf("file \"%s\" was copied to some of the destinations, but not all", context.filename))
	}

	fileCopier.stats.NumberOfSourceFiles++

	return
}

func (fileCopier *fileCopier) copyFile(context *copyContext, destinationPath string) (bool, error) {
	var err error

	// *** this is the main focus of this entire Go program, copying files to specific destinations ***

	sourceFilename := path.Join(fileCopier.config.source, context.subFolderPath, context.filename)
	destFilename := path.Join(destinationPath, context.filename)

	fileinfoSource, err := os.Stat(sourceFilename)
	if err != nil {
		return false, err
	} else if !fileinfoSource.Mode().IsRegular() {
		return false, fmt.Errorf("%s was not copied as it is not a regular file", sourceFilename)
	}

	// check to see if the file exists, and if it does,
	// then check the configuration to see if it should be replaced
	fileinfoDest, fileExists := fileCopier.doesDestFileExist(destFilename)
	if fileExists {
		if !fileCopier.checkIfFileShouldBeReplaced(context, fileinfoSource, fileinfoDest) {
			// the file should not be replaced
			return false, nil
		}
	}

	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		return false, err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destFilename)
	if err != nil {
		PrintError(fmt.Sprintf("error creating %s: %s", destFilename, err))
		return false, err
	}
	defer destFile.Close()

	bytesWritten, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return false, err
	}

	// flush file to storage and close it BEFORE changing the modified time of the file
	err = destFile.Sync()
	if err != nil {
		return false, err
	}
	destFile.Close()

	// update the access and modified time for the file to be that of the original file
	err = os.Chtimes(destFilename, fileinfoSource.ModTime(), fileinfoSource.ModTime())
	if err != nil {
		PrintError(fmt.Sprintf("failed to changed modified time: %s", err))
	}

	fileCopier.stats.TotalFilesCopied++
	fileCopier.stats.BytesCopied += bytesWritten

	return true, nil
}

func (fileCopier *fileCopier) doesDestFileExist(destFilename string) (os.FileInfo, bool) {
	var fileExists = false

	fileinfoDest, err := os.Stat(destFilename)
	if err == nil {
		fileExists = true
	} else if err != nil {
		if !os.IsNotExist(err) {
			fileExists = true
		}
	}

	return fileinfoDest, fileExists
}

func (fileCopier *fileCopier) checkIfFileShouldBeReplaced(context *copyContext, fileinfoSource os.FileInfo, fileinfoDest os.FileInfo) bool {
	returnValue := true

	switch fileCopier.config.replace {
	case replaceNever:
		fileCopier.stats.TotalFilesSkipped++
		warningMsg := fmt.Sprintf("%s was not copied to %s as it already exists, and the replace flag is set to \"never\"",
			context.filename, context.destinationPath)
		PrintWarning(warningMsg)
		returnValue = false
		break

	case replaceSkipIfSame:
		if (fileinfoSource.ModTime() == fileinfoDest.ModTime()) && (fileinfoSource.Size() == fileinfoDest.Size()) {
			fileCopier.stats.TotalFilesSkipped++
			warningMsg := fmt.Sprintf("%s was not copied to %s because it matches the datetime and size of an existing file, and the replace flag is set to \"skip\"",
				context.filename, context.destinationPath)
			PrintWarning(warningMsg)
			returnValue = false
		}
		break
	}

	return returnValue
}

func (fileCopier *fileCopier) handleFinish() {
	recovery := recover()
	if recovery != nil {
		PrintError(fmt.Sprintf("panic occurred:\n    %v", recovery))
	}
}
