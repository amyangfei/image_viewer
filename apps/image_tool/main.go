package main

import (
	"flag"
	"github.com/amyangfei/image_viewer/viewer"
	"log"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("Usage:\n  hello MOUNTPOINT URL")
	}
	viewer.Serve(flag.Arg(0), flag.Arg(1))
}
