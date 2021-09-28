package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asahnoln/go-scheduler/cli"
)

func main() {
	file, err := os.OpenFile(os.Args[1], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Welcome to the Scheduler!")
	c := cli.NewCLI(os.Stdin, os.Stdout)
	c.DB(file)
	log.Fatal(c.MainLoop())
}
