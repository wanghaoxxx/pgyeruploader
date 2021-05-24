# 一个简单的蒲公英上传工具


### 编译 

##### Windows
``
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o pgyeruploader main.go
``

##### Mac
``
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o pgyeruploader main.go
``

### 使用方式

``
上传 flutter release包 到蒲公英 
./pgyuploader -t release
``

``
上传 flutter debug 到蒲公英 ，默认就是debug包
./pgyuploader
``

