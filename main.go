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
		copylib.PrintError("the operation flag is required; it defines which operation in the config to execute...")
		os.Exit(2)
	}

	copyRunner := copylib.NewRunner(operation)
	go copyRunner.Copy()

	copyRunner.Waiter.Wait()

	stats := color.New(color.FgBlue, color.Bold)
	copylib.PrintColor(stats, "\nStats:")
	copylib.PrintStats("    Files Copied: ", fmt.Sprintf("%d", copyRunner.Stats.FilesCopied))
	printer := message.NewPrinter(language.English)
	copylib.PrintStats("    Bytes Copied: ", fmt.Sprintf("%d", copyRunner.Stats.FilesCopied))
	copylib.PrintStats("    Files Copied: ", printer.Sprintf("%d", copyRunner.Stats.BytesCopied))
	copylib.PrintStats("    Time to Copy: ", fmt.Sprintf("%f", copyRunner.Stats.TimeToCopy.Seconds()))
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
		copylib.PrintError(errOutput)
		copylib.PrintError("exiting...")

		os.Exit(1)
	}
}
