package main

import (
	"log"
	"os"
	"shop-reviews/cmd"
)

func main() {
	err := cmd.App.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
