package copylib

import (
	"github.com/fatih/color"
)

type LogMode int8

const (
	LogSilent LogMode = iota
	LogSimple
	LogDebug
	LogInfo
	LogVerbose
)

var currentLogMode = LogInfo

func SetLogMode(logMode LogMode) {
	// set the current logging mode
	currentLogMode = logMode
}

func Print(formattedString string) {
	if currentLogMode > LogSilent {
		color.Green("%s%s", color.GreenString("go-copy: "), color.WhiteString(formattedString))
	}
}

func PrintVersionInfo(stringOne string, stringTwo string) {
	color.Green("%s%s", color.CyanString(stringOne), color.MagentaString(stringTwo))
}

func PrintDebug(formattedString string) {
	if currentLogMode >= LogDebug {
		color.New(color.FgGreen, color.Bold).Print("go-copy: ")
		color.New(color.FgHiMagenta, color.Italic).Printf("%s\n", formattedString)
	}
}

func PrintInfo(formattedString string) {
	if currentLogMode >= LogInfo {
		color.Yellow("%s%s", color.WhiteString("go-copy: "), color.New(color.FgCyan, color.Italic).Sprint(formattedString))
	}
}

func PrintWarning(formattedString string) {
	if currentLogMode >= LogSimple {
		color.Yellow("%s%s", color.YellowString("go-copy: "), color.New(color.FgMagenta, color.Italic).Sprint(formattedString))
	}
}

func PrintError(formattedString string) {
	color.Red("%s%s", color.RedString("go-copy: "), color.New(color.FgBlue, color.Italic).Sprint(formattedString))
}

func PrintErrorHighlight(formattedString string) {
	color.Red("%s%s", color.RedString("go-copy: "), color.New(color.FgMagenta, color.Italic).Sprint(formattedString))
}

func PrintStats(stringOne string, stringTwo string) {
	color.Green("%s%s", color.GreenString(stringOne), color.MagentaString(stringTwo))
}

func PrintKeyValue(stringOne string, stringTwo string) {
	color.New(color.FgBlue, color.Bold).Printf("%s", stringOne)
	color.New(color.FgMagenta).Printf("%s\n", stringTwo)
}

func PrintKeyValueArray(stringOne string, stringArray []string) {
	color.New(color.FgBlue, color.Bold).Printf("%s\n", stringOne)
	for _, str := range stringArray {
		color.New(color.FgMagenta).Printf("    %s\n", str)
	}
}

func PrintColor(color *color.Color, formattedString string) {
	color.Println(formattedString)
}
