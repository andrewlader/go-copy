package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fatih/color"
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

	stats := color.New(color.FgBlue, color.Bold)
	stats.Println("\nStats:")
	color.Green("%s%s", color.GreenString("    Files Copied: "), color.MagentaString(fmt.Sprintf("%d", runner.Stats.FilesCopied)))
	printer := message.NewPrinter(language.English)
	color.Green("%s%s", color.GreenString("    Bytes Copied: "), color.MagentaString(printer.Sprintf("%d", runner.Stats.BytesCopied)))
	color.Green("%s%s", color.GreenString("    Time to Copy: "), color.MagentaString(fmt.Sprintf("%f", runner.Stats.TimeToCopy.Seconds())))
	color.White("\nAll done...\n\n")
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
