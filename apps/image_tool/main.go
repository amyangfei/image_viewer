package main

import (
	"fmt"
	"os"

	"github.com/amyangfei/image_viewer/viewer"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Headless bool `long:"headless" description:"whether use browser automation framework"`

	DriverPort int `long:"driver-port" default:"9515" description:"chrome driver port"`

	ShowVersion bool `long:"version" description:"print version"`

	FsInfo struct {
		MountPoint string
		Url        string
	} `positional-args:"yes" required:"yes" description:"mount point and crawling url"`
}

func main() {
	args := make([]string, len(os.Args)-1)
	copy(args, os.Args[1:])

	args, err := flags.NewParser(&opts, flags.PassDoubleDash|flags.HelpFlag).ParseArgs(args)
	if err != nil {
		if opts.ShowVersion {
			fmt.Println(viewer.Version("image_tool"))
			return
		}
		fmt.Println(err)
		return
	}

	if opts.ShowVersion {
		fmt.Println(viewer.Version("image_tool"))
		return
	}

	fsOpts := viewer.NewOptions()
	fsOpts.Headless = opts.Headless
	fsOpts.DriverPort = opts.DriverPort

	viewer.Serve(opts.FsInfo.MountPoint, opts.FsInfo.Url, fsOpts)
}
