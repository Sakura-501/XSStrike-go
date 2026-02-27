package main

import (
	"flag"
	"fmt"

	"github.com/Sakura-501/XSStrike-go/internal/ui"
	"github.com/Sakura-501/XSStrike-go/internal/version"
)

func main() {
	showVersion := flag.Bool("version", false, "show version")
	showVersionShort := flag.Bool("v", false, "show version")

	flag.Usage = func() {
		fmt.Print(ui.Banner())
		fmt.Println("Usage: xsstrike-go [options]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion || *showVersionShort {
		fmt.Printf("%s %s\n", version.AppName, version.Version)
		return
	}

	fmt.Print(ui.Banner())
	fmt.Println("Bootstrap complete. Next commits will add migrated XSStrike features.")
}
