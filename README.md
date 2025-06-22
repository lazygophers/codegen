# codegen

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/lazygophers/codegen)

> 根据proto文件生成代码，并且将相关标记写入 `pb.go` 中

这个项目是一个基于代码生成帮助快速开发 CURD 类的项目，你可以使用它来快速生成代码。

当前阶段，这个项目还在开发中，所以还有很多功能没有实现，如果你有兴趣，可以一起来完善这个项目。

帮助文档：https://lazygophers.pages.dev/

模板仓库：https://github.com/lazygophers/template

## 生成样例

首先，需要创建一个文件夹存放 `proto` 文件，假设文件夹名为 `proto`，并且在 `proto` 文件夹下存在一个 `hello.proto` 文件，内容如下

```proto
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

## 想要了解更多信息，可参考以下内容

- [DeepWiki](https://deepwiki.com/lazygophers/codegen)
