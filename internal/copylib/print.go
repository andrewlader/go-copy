package copylib

import (
	"github.com/fatih/color"
)

func Print(formattedString string) {
	if currentLogMode != LogSilent {
		color.Green("%s%s", color.GreenString("go-copy: "), color.WhiteString(formattedString))
	}
}

func PrintVersionInfo(stringOne string, stringTwo string) {
	color.Green("%s%s", color.CyanString(stringOne), color.MagentaString(stringTwo))
}

func PrintWarning(formattedString string) {
	if currentLogMode == LogVerbose {
		color.Yellow("%s%s", color.YellowString("go-copy: "), color.MagentaString(formattedString))
	}
}

func PrintError(formattedString string) {
	color.Red("%s%s", color.RedString("go-copy: "), color.CyanString(formattedString))
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
