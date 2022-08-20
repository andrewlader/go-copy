package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"gitlab.com/andrewlader/go-copy/copylib"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var operation string
var pauseAtEnd bool
var finishedSuccessfully bool

func init() {
	defer handleExit()

	finishedSuccessfully = false

	parseArguments()

	viper.SetConfigName("go-copy-config") // name of config file (without extension)
	viper.SetConfigType("yml")            // REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")              // path to look for the config file in
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		panic(fmt.Errorf("fatal error processing config file: %s", err))
	}
}

func main() {
	defer handleExit()

	if len(operation) < 1 {
		panic("the operation flag is required; it defines which operation in the config to execute...")
	}

	copyFileRunner := copylib.NewRunner(operation)
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
	color.White("\nAll done...\n\n")

	finishedSuccessfully = true

	if pauseAtEnd {
		pauseOutput()
	}
}

func parseArguments() {
	flag.StringVar(&operation, "operation", "", "defines the operation to execute (required)")
	flag.BoolVar(&pauseAtEnd, "pause", false, "determines if the app will pause before ending (optional)")

	flag.Parse()
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
