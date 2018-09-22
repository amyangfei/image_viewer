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
## Headless模式爬取

默认不执行js，如果需要比较完整的js执行，需要开启`--headless`选项，使用chrome headless渲染页面。Ubuntu/Debian 下安装chrome相关依赖的示例如下：

```bash
curl -sSL https://dl.google.com/linux/linux_signing_key.pub | apt-key add -
echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list
apt-get update && apt-get install -y google-chrome-stable
wget -N https://chromedriver.storage.googleapis.com/2.42/chromedriver_linux64.zip && unzip chromedriver_linux64.zip
mv -f chromedriver /usr/local/bin/chromedriver
```

## TODO

- [ ] Add test case
- [ ] Dependency management
- [ ] Javascript simulator, eg chrome headless
- [ ] Better filename with urlencode
- [ ] Image type pre detection, used for filename without extension and acceleration for dir list
- [ ] CI support
- [ ] Duplicate url optimization
