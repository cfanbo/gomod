package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cfanbo/gomod/core"
	"github.com/cfanbo/gomod/core/log"
)

var (
	modFile     string
	versionFlag bool
	helpFlag    bool
	debugFlag   bool
)

func init() {
	flag.BoolVar(&versionFlag, "v", false, "print the version")
	flag.BoolVar(&helpFlag, "h", false, "help info")
	flag.BoolVar(&debugFlag, "d", false, "debug mode")
}

func main() {
	flag.Parse()

	if versionFlag {
		fmt.Println(core.FullVersion())
		return
	}

	log.Init()
	if debugFlag {
		log.SetDebugLevel(true)
	}

	if helpFlag {
		helpText := `gomod is a tool for analysis go.mod depend module statistics.

Usage:
  gomod [Flags] [github.com/user/repo]

The Flags are:
  -v	print current version
  -h    compile packages and dependencies
  -d	debug mode`

		fmt.Println(helpText)
		return
	}

	var f string
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] != "-d" && os.Args[i] != "vv" {
				f = os.Args[i]
				break
			}
		}
	}

	// core
	core.Enter(f)
}
