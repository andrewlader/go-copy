package copier

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Runner struct {
	configName string
	config     *configuration
}

func NewRunner(configName string) *Runner {
	runner := &Runner{
		configName: configName,
	}

	config := getConfiguration(configName)

	runner.config = config

	return runner
}

func (runner *Runner) Copy() {
	defer handleFinish()

	runner.walkPath("")
}

func (runner *Runner) walkPath(subPath string) {
	var file os.FileInfo

	context := &context{
		subPath: subPath,
	}

	currentPath := path.Join(runner.config.source, context.subPath)
	log.Printf("walking path --> %s", currentPath)

	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		log.Printf("skipping path: %v", err)
	}

	for _, file = range files {
		if !file.IsDir() {
			// copy the file
			context.filename = file.Name()
			err := runner.copyFileToDestinations(context)
			if err != nil {
				log.Printf("failed to copy file \"%s\" due to error: %s", file.Name(), err)
			}
		} else {
			// walk the sub-folder path
			runner.walkPath(file.Name())
		}
	}
}

func (runner *Runner) copyFileToDestinations(context *context) error {
	var err error

	for _, destPath := range runner.config.destinations {
		context.destinationPath = destPath
		err = runner.copyFile(context)
	}

	return err
}

func (runner *Runner) copyFile(context *context) error {
	var err error

	sourceFilename := path.Join(runner.config.source, context.subPath, context.filename)
	destPath := path.Join(context.destinationPath, context.subPath)
	destFilename := path.Join(destPath, context.filename)

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		log.Printf("creating sub-path: %s", destPath)
		err = os.Mkdir(destPath, os.ModeDir)
		if err != nil {
			log.Printf("failed to create sub-path: %s", destPath)
			return err
		}
	}

	sourceFileStat, err := os.Stat(sourceFilename)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourceFilename)
	}

	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destFilename)
	if err != nil {
		log.Printf("error creating %s: %s", destFilename, err)
		return err
	}
	defer destFile.Close()

	log.Printf("copying \"%s\" to destination \"%s\"", context.filename, destPath)

	bytesWritten, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	log.Printf("copied %d bytes", bytesWritten)

	// flush file to storage
	err = destFile.Sync()
	if err != nil {
		return err
	}
	destFile.Close()

	// update the access and modified time for the file to be that of the original file
	err = os.Chtimes(destFilename, sourceFileStat.ModTime(), sourceFileStat.ModTime())
	if err != nil {
		log.Printf("failed to changed modified time: %s", err)
	} else {
		log.Printf("changed modified time to %v", sourceFileStat.ModTime())
	}

	return err
}

func handleFinish() {
	recovery := recover()
	if recovery != nil {
		log.Printf("panic occurred:\n    %v", recovery)
	}
}
