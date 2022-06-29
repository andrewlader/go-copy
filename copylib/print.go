package copylib

import (
	"github.com/fatih/color"
)

func Print(formattedString string) {
	color.Green("%s%s", color.GreenString("go-copy: "), color.WhiteString(formattedString))
}

func PrintWarning(formattedString string) {
	color.Yellow("%s%s", color.YellowString("go-copy: "), color.MagentaString(formattedString))
}

func PrintError(formattedString string) {
	color.Red("%s%s", color.RedString("go-copy: "), color.CyanString(formattedString))
}

func PrintStats(stringOne string, stringTwo string) {
	color.Green("%s%s", color.GreenString(stringOne), color.MagentaString(stringTwo))
}

func PrintColor(color *color.Color, formattedString string) {
	color.Println(formattedString)
}
