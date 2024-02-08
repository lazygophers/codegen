# codegen

> 根据proto文件生成代码，并且将相关标记写入 `pb.go` 中

## 如何使用

1. 安装 [protoc](https://github.com/protocolbuffers/protobuf/releases) 并配置环境变量
2. 安装 [protoc-gen-go](https:// github.com/golang/protobuf) 并配置环境变量
   ```bash
	go install github.com/golang/protobuf/protoc-gen-go
   ```
3. 安装 [codegen](https://github.com/lazygophers/codegen)
   ```bash
	go install github.com/lazygophers/codegen/cmd@latest
	```

4. 执行命令
   ```bash
	codegen gen pb -i <proto文件路径>
   ```

### 注意事项

1. 目前更多还是写死的配置，后续会逐渐完善
2. 生成的代码目前会按照如下目录结构生成
```
- proto
	- xxx.proto
- <go_package>
	- <package>.pb.go
```


## 生成样例

首先，需要创建一个文件夹存放 `proto` 文件，假设文件夹名为 `proto`，并且在 `proto` 文件夹下存在一个 `hello.proto` 文件，内容如下

```proto
syntax = "proto3";
syntax = "proto3";

package user;

option go_package = "github.com/lazygophers/user";

message User {
	// @gorm: primary_key
	// @role: admin:wr;user:r
	string id = 1;
}
```

然后执行命令 (可以在任意目录执行，只测试的相对路径的场景，没有测试绝对路径场景)

```bash
codegen gen pb -i proto/hello.proto
```

那么会在 `github.com/lazygophers/user` 目录下生成一个 `hello.pb.go`

```go
...其余内容...

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @gorm: primary_key
	// @role: admin:wr;user:r
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" gorm:"primary_key" role:"admin:wr;user:r"`
}

...其余内容...
```

此时的目录结构如下

```
- proto
	- hello.proto
- github.com
	- lazygophers
		- hello
			- hello.pb.go
```
