package main

import (
	"flag"
	"fmt"
	"github.com/amyangfei/image_viewer/viewer"
	"os"
)

var (
	flagSet = flag.NewFlagSet("image_tool", flag.ExitOnError)

	headless    = flagSet.Bool("headless", false, "whether use browser automation framework")
	showVersion = flagSet.Bool("version", false, "print version string")
)

func main() {
	flagSet.Parse(os.Args[1:])

	if *showVersion {
		fmt.Println(viewer.Version("image_tool"))
		return
	}

	flag.Parse()
	if len(flag.Args()) < 2 {
		fmt.Printf("Usage:\n image_tool MOUNTPOINT URL\n")
		return
	}

	opts := viewer.NewOptions()
	opts.Headless = *headless

	viewer.Serve(flag.Arg(0), flag.Arg(1), opts)
}
