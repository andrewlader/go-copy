package backup

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type copyContext struct {
	filename        string
	sourcePath      string
	destinationPath string
	subFolderPath   string
}

type stats struct {
	FilesCopied int
	BytesCopied int64
	TimeToCopy  time.Duration
}

type fileCopier struct {
	config *configuration
	stats  stats
}

func (fileCopier *fileCopier) run(config *configuration) {
	defer fileCopier.handleFinish()

	fileCopier.config = config

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
		log.Printf("skipping path %s: %v", currentPath, err)
	} else {
		for _, file = range files {
			if file.Mode().IsRegular() {
				// copy the file
				context.filename = file.Name()
				err := fileCopier.copyFileToDestinations(context)
				if err != nil {
					log.Printf("failed to copy file \"%s\" due to error: %s", file.Name(), err)
				}
			}
			if file.IsDir() {
				// walk the sub-folder path
				fileCopier.walkPath(path.Join(context.subFolderPath, file.Name()))
			}
		}
	}
}

func (fileCopier *fileCopier) copyFileToDestinations(context *copyContext) error {
	var err error

	log.Printf("copying file \"%s\"", context.filename)

	for _, destPath := range fileCopier.config.destinations {
		context.destinationPath = destPath

		// make sure all the sub folders exist for this destination path
		destinationPath := path.Join(context.destinationPath, context.subFolderPath)
		err = fileCopier.createSubFolders(context, destinationPath)
		if err != nil {
			continue
		}

		err = fileCopier.copyFile(context, destinationPath)
	}

	return err
}

func (fileCopier *fileCopier) copyFile(context *copyContext, destinationPath string) error {
	var err error

	// this is the main focus of this entire Go program, copying files to specific destinations

	sourceFilename := path.Join(fileCopier.config.source, context.subFolderPath, context.filename)
	destFilename := path.Join(destinationPath, context.filename)

	sourceFileStat, err := os.Stat(sourceFilename)
	if err != nil {
		return err
	} else if !sourceFileStat.Mode().IsRegular() {
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

	fileCopier.stats.FilesCopied++
	fileCopier.stats.BytesCopied += bytesWritten

	return nil
}

func (fileCopier *fileCopier) createSubFolders(context *copyContext, destinationPath string) error {
	var err error

	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		// some sub folders may exist already, so this will step down through each
		// sub folder and see if it needs to be created, and then move on to the
		// next child folder
		subPaths := strings.Split(context.subFolderPath, "/")
		newDir := context.destinationPath
		for _, subPath := range subPaths {
			newDir = path.Join(newDir, subPath)
			if _, err := os.Stat(newDir); os.IsNotExist(err) {
				err = os.Mkdir(newDir, os.ModeDir)
				if err != nil {
					log.Printf("failed to create sub-path: %s", newDir)
					break
				}
			}
		}
	}

	return err
}

func (fileCopier *fileCopier) handleFinish() {
	recovery := recover()
	if recovery != nil {
		log.Printf("panic occurred:\n    %v", recovery)
	}
}
