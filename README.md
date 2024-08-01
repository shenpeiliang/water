#### 使用说明

- 编译生成可执行文件：
```
go build -o water -ldflags "-s -w"
```

- 执行可执行文件（或在windows系统中双击water.exe运行程序）即可
```
./water
```

- 配置文件在config.toml，可修改水印文字等

- file目录放置需要加水印的文件，支持批量处理

- water目录放置生成的水印文件