package main

import (
	"fmt"

	"goweb/cmd/weble/web"
)

func main() {

	err := web.Run()
	fmt.Printf("Program exited with error: %s\n", err)
}
