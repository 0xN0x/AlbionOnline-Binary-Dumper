package main

import (
	"flag"
	"fmt"
	"os"

	dumper "github.com/0xN0x/albiononline-binary-dumper"
)

func main() {
	// Defining flags
	gameFolder := flag.String("g", "", "Game source folder")
	outputFolder := flag.String("o", "", "Output folder")

	// Parsing them
	flag.Parse()

	// Checking if required flags are set
	if *gameFolder == "" {
		flag.Usage()
		os.Exit(2)
	} else {
		if _, err := os.Stat(*gameFolder); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s does not exist\n", *gameFolder)
			os.Exit(2)
		}
	}

	if *outputFolder == "" {
		flag.Usage()
		os.Exit(2)
	} else {
		if _, err := os.Stat(*outputFolder); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s does not exist\n", *outputFolder)
			os.Exit(2)
		}
	}

	dumper.Dump(*gameFolder, *outputFolder)
}
