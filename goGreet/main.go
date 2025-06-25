package main

import (
	"fmt"

	gomodule "github.com/dvg1130/portfolio/goModule"
)

func main() {
	message := gomodule.Greet("Johnny")
	fmt.Println(message)
}
