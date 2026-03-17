package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/andrewlader/go-copy/internal/copylib"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var loadedConfigs bool
var version string
var date string
var commit string
var displayBuildInformation bool
var operation string
var listConfigs bool
var pauseAtEnd bool
var finishedSuccessfully bool
var logModeSilent bool
var logModeSimple bool
var logModeInfo bool
var logModeWarning bool
var logModeDebug bool
var logModeVerbose bool
var logMode copylib.LogMode

// init is called before the main function and is used to set up the configuration and handle any necessary initialization for the application.
func init() {
	defer handleExit()

	finishedSuccessfully = false

	parseArguments()

	if !displayBuildInformation {
		viper.SetConfigName("go-copy-config")         // name of config file (without extension)
		viper.SetConfigType("yml")                    // REQUIRED if the config file does not have the extension in the name
		viper.SetConfigType("yaml")                   // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(".")                      // look for the config file in this directory
		viper.AddConfigPath("/etc/go-copy")           // look for the config file in this directory
		viper.AddConfigPath("%USERPROFILE%/.go-copy") // look for the config file in this directory
		viper.AddConfigPath("%USERPROFILE%/.config")  // look for the config file in this directory
		viper.AddConfigPath("$HOME/")                 // look for the config file in this directory
		viper.AddConfigPath("$HOME/.config")          // look for the config file in this directory
		viper.AddConfigPath("%HOME/.config/go-copy")  // look for the config file in this directory
		viper.AddConfigPath("config/")                // look for the config file in this directory
		viper.AddConfigPath("configs/")               // look for the config file in this directory
		err := viper.ReadInConfig()                   // Find and read the config file
		if err != nil {
			loadedConfigs = false
			directUserToCreateConfigFile()
		} else {
			copylib.PrintInfo(fmt.Sprintf("config file loaded successfully: %s", viper.ConfigFileUsed()))
			loadedConfigs = true
		}
	}
}

// main is the entry point of the application. It handles command-line arguments and executes the appropriate actions based on those arguments.
func main() {
	defer handleExit()

	if displayBuildInformation {
		copylib.PrintVersionInfo("build version: ", version)
		copylib.PrintVersionInfo("build commit:  ", commit)
		copylib.PrintVersionInfo("build date:    ", date)
	} else if loadedConfigs {
		if listConfigs {
			copylib.ListConfigurations()
		} else {
			// run the main operation of the program, which is copying files based on the configuration
			finishedSuccessfully = runOperation()
		}
	}
}

// runOperation executes the file copy operation defined in the configuration.
func runOperation() bool {
	if len(operation) < 1 {
		panic("the operation flag is required; it defines which operation in the config to execute...")
	}

	copyFileRunner, err := copylib.NewRunner(operation)
	if err != nil {
		copylib.PrintError(fmt.Sprintf("error initializing runner for operation \"%s\": %s", operation, err))
		return false
	}

	go copyFileRunner.Copy()

	copyFileRunner.Waiter.Wait()

	stats := color.New(color.FgBlue, color.Bold)
	copylib.PrintColor(stats, "\nStats:")
	copylib.PrintStats("    Total Files Copied: ", fmt.Sprintf("%d (%d)", copyFileRunner.Stats.TotalFilesCopied/2, copyFileRunner.Stats.TotalFilesCopied))
	copylib.PrintStats("    Total Files Skipped: ", fmt.Sprintf("%d (%d)", copyFileRunner.Stats.TotalFilesSkipped/2, copyFileRunner.Stats.TotalFilesSkipped))
	copylib.PrintStats("    Number of Source Files: ", fmt.Sprintf("%d", copyFileRunner.Stats.NumberOfSourceFiles))
	copylib.PrintStats("    Number of Destinations: ", fmt.Sprintf("%d", copyFileRunner.Stats.NumberOfDestinations))
	printer := message.NewPrinter(language.English)
	copylib.PrintStats("    Bytes Copied: ", printer.Sprintf("%d", copyFileRunner.Stats.BytesCopied))
	copylib.PrintStats("    Time to Copy: ", fmt.Sprintf("%f", copyFileRunner.Stats.TimeToCopy.Seconds()))
	copylib.PrintStats("    Warnings: ", fmt.Sprintf("%d", copyFileRunner.Stats.NumberOfWarnings))
	copylib.PrintStats("    Errors: ", fmt.Sprintf("%d", copyFileRunner.Stats.NumberOfErrors))
	copylib.PrintStats("    Operation: ", operation)
	color.White("\nAll done...\n\n")

	if pauseAtEnd {
		pauseOutput()
	}

	return copyFileRunner.Stats.NumberOfErrors == 0
}

// directUserToCreateConfigFile prompts the user to create an empty YAML config file in the appropriate location for the OS.
func directUserToCreateConfigFile() {
	var defaultPath string
	var newFile string

	copylib.PrintError("No config file found, and this application lacks the permissions to create one...")
	switch runtimeOS := runtime.GOOS; runtimeOS {
	case "windows":
		userProfile := os.Getenv("USERPROFILE")
		if userProfile != "" {
			defaultPath = userProfile + "\\.go-copy"
		} else {
			defaultPath = ".\\.config"
		}
		newFile = defaultPath + "\\go-copy-config.yaml"
	case "darwin", "linux":
		defaultPath = "/etc/go-copy/"
		newFile = defaultPath + "go-copy-config.yaml"
	default:
		defaultPath = "./configs/"
		newFile = defaultPath + "go-copy-config.yaml"
	}

	copylib.PrintError("It is recommended to create an empty config file here:")
	copylib.PrintErrorHighlight(fmt.Sprintf("    %s", newFile))
}

// parseArguments processes the command-line arguments and sets the appropriate variables.
func parseArguments() {
	flag.BoolVar(&displayBuildInformation, "version", false, "display build & version information")
	flag.StringVar(&operation, "operation", "", "defines the operation to execute (required)")
	flag.BoolVar(&listConfigs, "list", false, "list all backup sets in the config")
	flag.BoolVar(&pauseAtEnd, "pause", false, "determines if the app will pause before ending (optional)")
	flag.BoolVar(&logModeSilent, "silent", false, "logging out put will be sparse (optional)")
	flag.BoolVar(&logModeSimple, "simple", false, "logging out put will be normal (optional)")
	flag.BoolVar(&logModeWarning, "warning", false, "logging out put will be at the warning level (optional)")
	flag.BoolVar(&logModeInfo, "info", false, "logging out put will be at the info level (optional)")
	flag.BoolVar(&logModeDebug, "debug", false, "logging out put will be at the debug level (optional)")
	flag.BoolVar(&logModeVerbose, "verbose", false, "logging out put will be verbose (optional)")

	flag.Parse()

	if logModeSilent {
		logMode = copylib.LogSilent
	} else if logModeSimple {
		logMode = copylib.LogSimple
	} else if logModeWarning {
		logMode = copylib.LogWarning
	} else if logModeInfo {
		logMode = copylib.LogInfo
	} else if logModeDebug {
		logMode = copylib.LogDebug
	} else if logModeVerbose {
		logMode = copylib.LogVerbose
	} else {
		logMode = copylib.LogInfo
	}

	copylib.SetLogMode(logMode)
}

// pauseOutput prompts the user to press enter before continuing, effectively pausing the output.
func pauseOutput() {
	copylib.PrintAlways("Press enter to continue...")
	fmt.Scanln()
}

// handleExit recovers from any panics that occur during the execution of the program and prints an error message before exiting.
// If the program finishes successfully, it prints a success message.
func handleExit() {
	recovery := recover()
	if recovery != nil {
		errOutput := fmt.Sprintf("panic occurred:\n    %v", recovery)
		copylib.PrintError(errOutput)
		copylib.PrintError("go-copy has stopped with an error")
		copylib.PrintError(fmt.Sprintf("%s", debug.Stack()))
		os.Exit(1)
	} else if finishedSuccessfully {
		if len(operation) > 0 {
			copylib.PrintAlways(fmt.Sprintf("go-copy has completed operation \"%s\" successfully", operation))
		}
	}
}
