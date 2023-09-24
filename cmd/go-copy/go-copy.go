package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"gitlab.com/andrewlader/go-copy/internal/copylib"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:embed git-describe.txt
var buildInfo string

var displayBuildInformation bool
var operation string
var pauseAtEnd bool
var finishedSuccessfully bool
var logModeSilent bool
var logModeSimple bool
var logModeVerbose bool
var logMode copylib.LogMode

func init() {
	defer handleExit()

	finishedSuccessfully = false

	parseArguments()

	viper.SetConfigName("go-copy-config")           // name of config file (without extension)
	viper.SetConfigType("yml")                      // REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType("yaml")                     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                        // path to look for the config file in
	viper.AddConfigPath("cmd/")                     // path to look for the config file in
	viper.AddConfigPath("config/")                  // path to look for the config file in
	viper.AddConfigPath("configs/")                 // path to look for the config file in
	err := viper.ReadInConfig()                     // Find and read the config file
	if (!displayBuildInformation) && (err != nil) { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error processing config file: %s", err))
	}
}

func main() {
	defer handleExit()

	if displayBuildInformation {
		buildInformation := strings.Split(buildInfo, "\n")
		if len(buildInformation) > 0 {
			copylib.PrintVersionInfo("go-copy version: ", buildInformation[0])
		}
		if len(buildInformation) > 1 {
			copylib.PrintVersionInfo("go version:      ", buildInformation[1])
		}
		if len(buildInformation) > 2 {
			copylib.PrintVersionInfo("build date:      ", buildInformation[2])
		}
	} else {
		if len(operation) < 1 {
			panic("the operation flag is required; it defines which operation in the config to execute...")
		}

		copyFileRunner := copylib.NewRunner(operation, logMode)
		go copyFileRunner.Copy()

		copyFileRunner.Waiter.Wait()

		stats := color.New(color.FgBlue, color.Bold)
		copylib.PrintColor(stats, "\nStats:")
		copylib.PrintStats("    Number of Source Files: ", fmt.Sprintf("%d", copyFileRunner.Stats.NumberOfSourceFiles))
		copylib.PrintStats("    Number of Destinations: ", fmt.Sprintf("%d", copyFileRunner.Stats.NumberOfDestinations))
		copylib.PrintStats("    Total Files Skipped: ", fmt.Sprintf("%d", copyFileRunner.Stats.TotalFilesSkipped))
		copylib.PrintStats("    Total Files Copied: ", fmt.Sprintf("%d", copyFileRunner.Stats.TotalFilesCopied))
		printer := message.NewPrinter(language.English)
		copylib.PrintStats("    Bytes Copied: ", printer.Sprintf("%d", copyFileRunner.Stats.BytesCopied))
		copylib.PrintStats("    Time to Copy: ", fmt.Sprintf("%f", copyFileRunner.Stats.TimeToCopy.Seconds()))
		color.Cyan("\nCopy complete for operation: ", operation)
		color.White("\nAll done...\n\n")

		finishedSuccessfully = true

		if pauseAtEnd {
			pauseOutput()
		}
	}
}

func parseArguments() {
	flag.BoolVar(&displayBuildInformation, "version", false, "display build & version information")
	flag.StringVar(&operation, "operation", "", "defines the operation to execute (required)")
	flag.BoolVar(&pauseAtEnd, "pause", false, "determines if the app will pause before ending (optional)")
	flag.BoolVar(&logModeSilent, "silent", false, "logging out put will be sparse (optional)")
	flag.BoolVar(&logModeSimple, "simple", false, "logging out put will be normal (optional)")
	flag.BoolVar(&logModeVerbose, "verbose", false, "logging out put will be verbose (optional)")

	flag.Parse()

	if logModeSilent {
		logMode = copylib.LogSilent
	} else if logModeSimple {
		logMode = copylib.LogSimple
	} else if logModeVerbose {
		logMode = copylib.LogVerbose
	}
}

func pauseOutput() {
	copylib.Print("Press enter to continue...")
	fmt.Scanln()
}

func handleExit() {
	recovery := recover()
	if recovery != nil {
		errOutput := fmt.Sprintf("panic occurred:\n    %v", recovery)
		copylib.PrintError(errOutput)
		copylib.PrintError("go-copy has stopped with an error")

		os.Exit(1)
	} else if finishedSuccessfully {
		copylib.PrintError("go-copy has completed its job successfully")
	}
}
