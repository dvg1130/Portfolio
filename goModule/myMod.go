package gomodule

import "fmt"

func Greet(name string) string {
	greeting := fmt.Sprintf("Hi, %v. Welcome", name)
	return greeting

}
