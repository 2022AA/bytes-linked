##byteds-linked后端部署

#####1. 项目根目录下

#####2. 编译，生成二进制文件
```
GOOS=linux GOARCH=amd64 go build -o bytes-linked main.go
```
#####3. 拷贝二进制文件和conf文件到服务器相应目录

#####4. 配置文件调整相应的参数（参考服务器上）

#####5. 后台执行
```
nohup ./bytes-linked > nohup.out &
```

#####注：本地运行需要修改配置文件

```
TempLocalRootDir = ["你bytes-linked实际的目录"/bytes-linked/tmp/]
```