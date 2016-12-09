package main

import (
	"github.com/mattyw/mtty"
	"os"
)

func main() {
	mtty.Loop(os.Stdin, os.Stdout, os.Stderr)
}
