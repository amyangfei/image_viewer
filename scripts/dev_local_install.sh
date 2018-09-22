#!/bin/bash

# install dependcies for local developing

cur=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
src_dst=$GOPATH/src/github.com/amyangfei/image_viewer

rm -rf $src_dst
mkdir -p $src_dst
cp -r $cur/../viewer $src_dst

cd $src_dst/viewer
echo "installing image_viewer/viewer"
go install

echo "done!"
