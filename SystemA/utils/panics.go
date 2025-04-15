package utils

import "fmt"

func PanicIfError(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("\n\n %s \nError %s\n\n", RedText(message), YellowText(err.Error())))
	}
}
