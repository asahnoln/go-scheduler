package main

import (
	"fmt"
	"os"

	"github.com/asahnoln/go-scheduler"
)

func main() {
	fmt.Println("Welcome to the Scheduler!")
	c := scheduler.NewCLI(os.Stdin, os.Stdout)
	for {
		err := c.Process()

		if err != nil {
			panic(err)
		}
	}
}
