package backup

import "fmt"

func printF(formattedString string, args ...string) {
	fmt.Printf(formattedString, args)
	fmt.Println()
}

func printErrorF(formattedString string, args ...string) {
	fmt.Printf(formattedString, args)
	fmt.Println()
}
