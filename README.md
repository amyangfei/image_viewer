# Simple Image Crawling and File System Mapping

Have fun

## Build

```bash
export build_path=/path/to/build
mkdir -p $build_path/src/github.com/amyangfei && cd $build_path/src/github.com/amyangfei
git clone https://github.com/amyangfei/image_viewer
export GOPATH=$GOPATH:$build_path
cd image_viewer && make
```
## Headless Crawling

Javascript executing is turned off by default. If we want to execute js, turn on `--headless` option and chrome headless will be used. Chrome and chrome driver is needed in headless mode. Dependencies installation instructions in Ubuntu/Debian is following:

```bash
curl -sSL https://dl.google.com/linux/linux_signing_key.pub | apt-key add -
echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list
apt-get update && apt-get install -y google-chrome-stable
wget -N https://chromedriver.storage.googleapis.com/2.42/chromedriver_linux64.zip && unzip chromedriver_linux64.zip
mv -f chromedriver /usr/local/bin/chromedriver
```

## TODO

- [ ] Add test case
- [x] Dependency management
- [x] Javascript simulator, eg chrome headless
- [ ] Better filename against urlencode
- [ ] Image type pre detection, used for filename without extension and acceleration for dir list
- [ ] CI support
- [ ] Duplicate url optimization
- [ ] Better url and img src extract strategy
