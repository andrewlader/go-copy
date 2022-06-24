package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"gitlab.com/andrewlader/go-copy/backup"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var operation string

func init() {
	defer handleExit()

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
		backup.PrintError("the operation flag is required; it defines which operation in the config to execute...")
		os.Exit(2)
	}

	runner := backup.NewRunner(operation)
	go runner.Copy()

	runner.Waiter.Wait()

	stats := color.New(color.FgBlue, color.Bold)
	backup.PrintColor(stats, "\nStats:")
	backup.PrintStats("    Files Copied: ", fmt.Sprintf("%d", runner.Stats.FilesCopied))
	printer := message.NewPrinter(language.English)
	backup.PrintStats("    Bytes Copied: ", fmt.Sprintf("%d", runner.Stats.FilesCopied))
	backup.PrintStats("    Files Copied: ", printer.Sprintf("%d", runner.Stats.BytesCopied))
	backup.PrintStats("    Time to Copy: ", fmt.Sprintf("%f", runner.Stats.TimeToCopy.Seconds()))
	color.White("\nAll done...\n\n")
}

func parseArguments() {
	flag.StringVar(&operation, "operation", "", "defines the operation to execute (required)")

	flag.Parse()
}

func handleExit() {
	recovery := recover()
	if recovery != nil {
		errOutput := fmt.Sprintf("panic occurred:\n    %v", recovery)
		backup.PrintError(errOutput)
		backup.PrintError("exiting...")

		os.Exit(1)
	}
}
