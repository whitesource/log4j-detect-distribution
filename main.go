package main

import (
	"github.com/whitesource/log4j-detect/cmd"
	"log"
)

func main() {
	if err := cmd.NewCmdRoot().Execute(); err != nil {
		log.Fatalln(err)
	}
}
