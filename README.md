# codegen

> 根据proto文件生成代码，并且将相关标记写入 `pb.go` 中

这个项目是一个基于代码生成帮助快速开发 CURD 类的项目，你可以使用它来快速生成代码。

当前阶段，这个项目还在开发中，所以还有很多功能没有实现，如果你有兴趣，可以一起来完善这个项目。

## 使用前必读

在前期阶段，我们需要一些基本共识：

- 相关工具类将会优先使用本 Orgaization 下的工具
- 代码风格将会优先使用本 Orgaization 下的代码风格
- 生成是基于 protobuf 的，所以你需要了解 protobuf 的基本使用，并且在在你的项目中使用 protobuf 进行接口、接口类型、接口文档的定义
- 为了简化实现，不区分空值（默认值）与未传入的区别，所以在使用时需要注意对空值的处理
- 所有的数据表的类型名都需要以 `Model` 为开头，例如 `ModelUser`，`ModelOrder` 等
	- 除非用户主动开启，否则不会自动执行以下操作
		- 数据表默认使用 `id` 作为主键，如果存在字段名为 `id` 的字段，那么将会使用该字段作为主键，并且主动为其添加
		  相关 `gorm` 标签
		- 数据表默认使用 `created_at`、`updated_at`、`deleted_at`
		  作为时间字段，如果存在字段名为 `created_at`、`updated_at`、`deleted_at` 的字段，那么将会使用该字段作为时间字段，并且主动为其添加
		  相关 `gorm` 标签
- 所有接口的请求、返回类型都需要为 `Req`、`Resp` 结尾，例如 `UserListReq`、`UserListResp` 等
- 常用CURD的命名规则为
	- List: `ListUserReq`、`ListUserResp`
		- List的分页方式先才用offset、limit，后续再考虑更多的分页方式例如cursor、page或者token等方式
		- List的筛选、排序方式优先采用fliter(key:value)的方式，后续再考虑更多的筛选、排序方式
		- Get（依据Id或唯一标识字段进行获取数据）: `GetUserReq`、`GetUserResp`
		-
		GetBy（依据非唯一标识字段进行获取数据或根据非唯一字符按照一定规则获取）: `GetUserByUsernameReq`、`GetUserByUsernameResp`
		- Create: `AddUserReq`、`AddUserResp`
		- Update: `UpdateUserReq`、`UpdateUserResp`
		- Delete: `DelUserReq`、`DelUserResp`
		- Set: `SetUserReq`、`SetUserResp`
- 整体项目结构
	- xxx.pb.go
	- xxx.proto(单体项目推荐放在当前目录，微服务推荐放在公共仓库)
	- go.mod
	- go.sum
	- .editorconfig
    - orm_gen.go(gorm 配套的接口协议的生成)
    - table_gen.go(表名的生成)
	- state
		- table.go(数据库相关的内容)
        - cache.go(缓存相关的内容)
        - config.go(配置文件相关内容)
        - state.go(有状态的服务配置的初始化)

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
	codegen gen pb <action> -i <proto文件路径>
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

3. 如果全局变量无法找到，可以尝试在配置文件中使用绝对路径

### 配置文件

> 配置文件名 `codegen.cfg.yaml`

会依次在环境变量（LAZYGO_CODEGEN_CONFIG_FILE）、用户配置路径（/root/.config/lazygophers）依次查找，如果没有找到则使用默认配置

```yaml
protoc_path: "<protoc的绝对路径>"
protogen_go_path: "<protoc-gen-go的绝对路径>"

# 生成的代码的目录，即 protoc 执行的路径，实际生成的代码将会放在 <output_path>/<go_package> 下
# 如果不指定，默认规则为 <proto_file_path>/../<go_package>
output_path: "<生成的代码的目录>"
go_module_prefix: "<go module的前缀，在生成时会拼接proto中的 go_package 字段当做包名>"

# 自定义模版文件，要求 .gtpl 格式
template:
	# .editorconfig 模板文件路径，用于 IDE 的默认格式化
	editorconfig: "<.editorconfig 文件路径，如果不存在则使用默认>"
	# orm.go 模板文件路径，用于 gorm 的相关数据生成
	orm: "<orm.go 模板文件路径，如果不存在则是使用默认>"
	# table_name.go 模板文件路径，用于 gorm 的相关表名
	table_name: "<table_name.go 模板文件路径，如果不存在则是使用默认>"

	proto_files:
		- ~/proto/
		- ~/core/

	# .proto 文件相关的模板文件
	proto:
		# rpc 服务的模板文件路径
		service: "proto.service 的模板文件路径"

		# rpc 相关的
		rpc:
			# key → action value → 模版文件路径
			add:
				name: "rpc 名字的生成模板"
				req: "message.request 生成模板"
				resp: "messsage.response 生成模板"
			set:
				name: "rpc 名字的生成模板"
				req: "message.request 生成模板"
				resp: "messsage.response 生成模板"
			update:
				name: "rpc 名字的生成模板"
				req: "message.request 生成模板"
				resp: "messsage.response 生成模板"
			del:
				name: "rpc 名字的生成模板"
				req: "message.request 生成模板"
				resp: "messsage.response 生成模板"
			# 如果 key 为 list，会触发筛选项的生辰
			list:
				name: "rpc 名字的生成模板"
				req: "message.request 生成模板"
				resp: "messsage.response 生成模板"

	# state/table.go 模板文件路径
	table: "<state/table.go 模板文件路径，如果不存在则是使用默认>"
	# state/conf.go 模板文件路径
	conf: "<state/conf.go 模板文件路径，如果不存在则是使用默认>"
	# state/cache.go 模板文件路径
	cache: "<state/cache.go 模板文件路径，如果不存在则是使用默认>"
	# state/state.go 模板文件路径
	state: "<state/state.go 模板文件路径，如果不存在则是使用默认>"

# 默认的标签，会追加到原始标签中，key 为标签名, value.key 为字段名, value.value 为标签值
default_tag:
	gorm:
		id: "column:id;primaryKey;autoIncrement;not null"
		created_at: "column:created_at;<-:create;autoCreateTime;not null"
		updated_at: "column:updated_at;<-;autoUpdateTime;not null"
		deleted_at: "column:deleted_at;autoDeleteTime"
		name: "type:varchar(255);not null"
		username: "type:varchar(255);not null"
		email: "type:varchar(255);not null"
	validate:
		name: "required,max=200,min=2"
		username: "required,max=100"
		email: "required,email"

# 数据表相关的配置
tables:
	# 是否关闭，默认为为 false，如果单独使用命令，则不会生效
	disable: false

	# 是否关闭自动生成的 gorm tag，默认为 false
	# 关闭自动生成 id 相关的gorm tag，默认为 false
	disable_field_id: false
	# 关闭自动生成 created_at 的gorm tag，默认为 false
	disable_field_created_at: false
	# 关闭自动生成 updated_at 的gorm tag，默认为 false
	disable_field_updated_at: false
	# 关闭自动生成 deleted_at 的gorm tag，默认为 false
	disable_field_deleted_at: false
	# 关闭自动生成 column 的gorm tag，默认为 false
	disable_gorm_tag_column: false
```

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
