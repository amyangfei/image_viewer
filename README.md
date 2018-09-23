# Simple Image Crawling and File System Mapping

[![Build Status](https://travis-ci.org/amyangfei/image_viewer.svg?branch=master)](https://travis-ci.org/amyangfei/image_viewer)

Have fun

## Introduction

This is a simple tool mapping images and sub links in a single html page to file system directory structure.

When we run `image_tool` a simple file system server will be running background. The file system server is based on `PathFileSystem` provided by [go-fuse](https://github.com/hanwen/go-fuse). File system operation such as `ls`, `cd`, `cat` will trigger interface defined in `go-fuse`, so we implement some useful interface in order to update file system structure dynamicly. Currently the file system information including dir entry list, file attributes and file data is all stored in memory.

## Build

```bash
$ export build_path=/path/to/build
$ mkdir -p $build_path/src/github.com/amyangfei && cd $build_path/src/github.com/amyangfei
$ git clone https://github.com/amyangfei/image_viewer
$ export GOPATH=$GOPATH:$build_path
$ cd image_viewer && make
```
## Headless Crawling

Javascript executing is turned off by default. If we want to execute js, turn on `--headless` option and chrome headless will be used. Chrome and chrome driver is needed in headless mode. Dependencies installation instructions in Ubuntu/Debian is following:

```bash
$ curl -sSL https://dl.google.com/linux/linux_signing_key.pub | apt-key add -
$ echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list
$ apt-get update && apt-get install -y google-chrome-stable
$ wget -N https://chromedriver.storage.googleapis.com/2.42/chromedriver_linux64.zip && unzip chromedriver_linux64.zip
$ mv -f chromedriver /usr/local/bin/chromedriver
```

## TODO

- [ ] Add test case
- [x] Dependency management
- [x] Javascript simulator, eg chrome headless
- [ ] Better filename against urlencode
- [ ] Image type pre detection, used for filename without extension and acceleration for dir list
- [x] CI support
- [x] Duplicate url optimization
- [ ] Better url and img src extract strategy
