package backup

import (
	"github.com/fatih/color"
)

func printF(formattedString string) {
	color.Green("%s%s", color.GreenString("go-copy: "), color.WhiteString(formattedString))
}

func printWarningF(formattedString string) {
	color.Yellow("%s%s", color.YellowString("go-copy: "), color.MagentaString(formattedString))
}

func printErrorF(formattedString string) {
	color.Red("%s%s", color.RedString("go-copy: "), color.CyanString(formattedString))
}
