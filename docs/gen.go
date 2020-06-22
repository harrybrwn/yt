package main

import (
	"fmt"
	"log"
	"os"

	"github.com/harrybrwn/yt/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	root := cmd.RootCommand()
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [directory]\n", os.Args[0])
		os.Exit(1)
	}
	if err := doc.GenManTree(root, nil, os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
