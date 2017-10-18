package main

import (
	"fmt"
	"os"
)

var Version string

func main() {
	for i, arg := range os.Args {
		switch arg {
		case "-h":
			fallthrough
		case "--help":
			fmt.Println("no one can help you")
			return
		case "-v":
			fallthrough
		case "--version":
			fmt.Println("ogrt-tool", Version)
			return
		case "--show-signature":
			f, err := os.Open(os.Args[i+1])
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			sigs, err := FindSignatures(f)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Found", len(sigs), "signatures:")
			for _, s := range sigs {
				fmt.Println(s)
			}
		case "--generate-load":
			generateLoad()
		}
	}
}
