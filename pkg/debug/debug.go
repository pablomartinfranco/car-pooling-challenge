package debug

import "fmt"

func Dump[T any](object T) {
	fmt.Printf("%+v\n", object)
}
