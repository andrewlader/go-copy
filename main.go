package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gitlab.com/andrewlader/go-copy/backup"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var operation string

func init() {
	parseArguments()

	viper.SetConfigName("go-copy-config") // name of config file (without extension)
	viper.SetConfigType("yml")            // REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")              // path to look for the config file in
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func main() {
	defer handleExit()

	if len(operation) < 1 {
		log.Fatal("The operation field is required. Exiting...")
	}

	runner := backup.NewRunner(operation)
	go runner.Copy()

	runner.Waiter.Wait()

	fmt.Print("Stats:\n")
	fmt.Printf("    Files Copied: %d\n", runner.Stats.FilesCopied)
	printer := message.NewPrinter(language.English)
	fmt.Print(printer.Sprintf("    Bytes Copied: %d\n", runner.Stats.BytesCopied))
	fmt.Printf("    Time to Copy: %f seconds\n", runner.Stats.TimeToCopy.Seconds())
}

func parseArguments() {
	flag.StringVar(&operation, "operation", "", "defines the operation to execute (required)")

	flag.Parse()
}

func handleExit() {
	recovery := recover()
	if recovery != nil {
		log.Printf("panic occurred:\n    %v", recovery)
		log.Print("exiting program")
	}
}
