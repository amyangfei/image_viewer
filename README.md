# simple image viewer

have fun

## Build

```bash
export build_path=/path/to/build
mkdir -p $build_path/src/github.com/amyangfei && cd $build_path/src/github.com/amyangfei
git clone https://github.com/amyangfei/image_viewer
export GOPATH=$GOPATH:$build_path
cd image_viewer && make
```

## TODO

- [ ] Add test case
- [ ] Dependency management
- [ ] Javascript simulator, eg chrome headless
- [ ] Better filename with urlencode
- [ ] Image type pre detection, used for filename without extension and acceleration for dir list
- [ ] CI support
- [ ] Duplicate url optimization
