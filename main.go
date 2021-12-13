//go:generate pkger
package main

import (
	"github.com/tyrm/supreme-robot/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err.Error())
	}
}
