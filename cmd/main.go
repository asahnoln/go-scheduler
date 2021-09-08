package main

import (
	"fmt"
	"os"

	"github.com/asahnoln/go-scheduler/cli"
)

func main() {
	fmt.Println("Welcome to the Scheduler!")
	c := cli.NewCLI(os.Stdin)
	w, err := c.Process()

	if err != nil {
		panic(err)
	}

	fmt.Println(w)
}
