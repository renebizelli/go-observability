package utils

var colorReset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var magenta = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

func YellowText(text string) string {
	return yellow + text + colorReset
}

func RedText(text string) string {
	return red + text + colorReset
}

func GreenText(text string) string {
	return green + text + colorReset
}

func BlueText(text string) string {
	return blue + text + colorReset
}

func MagentaText(text string) string {
	return magenta + text + colorReset
}

func CyanText(text string) string {
	return cyan + text + colorReset
}
