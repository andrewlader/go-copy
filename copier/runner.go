package copier

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const lineLength = 80

type Runner struct {
	Waiter     sync.WaitGroup
	configName string
	config     *configuration
	Stats      struct {
		FilesCopied int
		BytesCopied int64
		TimeToCopy  time.Duration
	}
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

	startTime := time.Now()

	runner.walkPath("")

	log.Printf("\r%-*s", lineLength, "file copy complete...")
	log.Println()

	runner.Stats.TimeToCopy = time.Now().Sub(startTime)
}

func (runner *Runner) walkPath(subPath string) {
	var file os.FileInfo

	context := &context{
		subPath: subPath,
	}

	currentPath := path.Join(runner.config.source, context.subPath)

	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		log.Printf("skipping path %s: %v", currentPath, err)
	} else {
		for _, file = range files {
			if file.Mode().IsRegular() {
				// copy the file
				context.filename = file.Name()
				err := runner.copyFileToDestinations(context)
				if err != nil {
					log.Printf("failed to copy file \"%s\" due to error: %s", file.Name(), err)
				}
			}
			if file.IsDir() {
				// walk the sub-folder path
				runner.walkPath(path.Join(context.subPath, file.Name()))
			}
		}
	}
}

func (runner *Runner) copyFileToDestinations(context *context) error {
	var err error

	log.Printf("copying file \"%s\"", context.filename)

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
		subPaths := strings.Split(context.subPath, "/")
		newDir := context.destinationPath
		for _, subPath := range subPaths {
			newDir = path.Join(newDir, subPath)
			if _, err := os.Stat(newDir); os.IsNotExist(err) {
				err = os.Mkdir(newDir, os.ModeDir)
				if err != nil {
					log.Printf("failed to create sub-path: %s", newDir)
					return err
				}
			}
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

	bytesWritten, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// flush file to storage and close it BEFORE changing the modified time of the file
	err = destFile.Sync()
	if err != nil {
		return err
	}
	destFile.Close()

	// update the access and modified time for the file to be that of the original file
	err = os.Chtimes(destFilename, sourceFileStat.ModTime(), sourceFileStat.ModTime())
	if err != nil {
		log.Printf("failed to changed modified time: %s", err)
	}

	runner.Stats.FilesCopied++
	runner.Stats.BytesCopied += bytesWritten

	return err
}

func (runner *Runner) handleFinish() {
	recovery := recover()
	if recovery != nil {
		log.Printf("panic occurred:\n    %v", recovery)
	}

	runner.Waiter.Done()
}
