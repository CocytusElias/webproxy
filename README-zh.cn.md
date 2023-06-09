# WebProxy

![standard-readme compliant](https://static.elias.ink/resource/202306092335092.svg)
![standard-readme compliant](https://static.elias.ink/resource/202306092335497.svg)
![standard-readme compliant](https://static.elias.ink/resource/202306092335909.svg)
![standard-readme compliant](https://static.elias.ink/resource/202306092335041.svg)
![standard-readme compliant](https://static.elias.ink/resource/202306092335423.svg)



这是一个仅支持 `http` 的轻量内网穿透工具，它并不比 [Ngrok](https://github.com/inconshreveable/ngrok)、[Frp](https://github.com/fatedier/frp) 更加强大。如果你需要更加强大的内网穿透代理工具，建议选用其他工具。



## 目录

- [WebProxy 是什么？](#WebProxy 是什么？)
- [运行](#运行)
- [运行](#运行)
- [配置](#配置)
- [模块](#模块)
- [维护者](#维护者)



## WebProxy 是什么？

`WebProxy` 是一个轻量内网 `web` 代理项目，目前仅支持 `HTTP` 协议的转发。

它并不比 [Ngrok](https://github.com/inconshreveable/ngrok)、[Frp](https://github.com/fatedier/frp) 更加强大。如果你需要更加强大的内网穿透代理工具，建议选用其他工具。



它与一般的代理服务一样，由 `Service` 与 `Client` 两部分组成。 `Service` 部署在公网服务器上（或是其他你需要的服务器上）， `Client` 则部署在能同时访问到内网服务和  `Service` 所在的服务器的机器上。



`Service` 与 `Client` 之间由 `websocket` 来进行通信。



整体的请求代理流程如下：

1.  `Service` 启动 `Http` 服务器，用来接受外部网络的 `Http` 请求，以及作为 `websocket` 接受 `Client` 的   `websocket` 请求。
2.  `Client` 向  `Service`  发起  `websocket` 请求，与 `Service` 建立起  `websocket` 链接。
3.  `Service` 来接受外部网络的 `Http` 请求，并将请求信息通过  `websocket` 转发给  `Client` 。
4.  `Client` 拿到请求后，根据配置来获取具体转发的服务地址，以及重写后的请求，并发起  `Http` 请求。
5.  `Client` 发起请求后，将响应内容通过  `websocket` 转发给 `Service`。
6.  `Service` 根据之前的记录，向发起请求的客户端返回响应内容。



## 运行

### 本地 Make 编译

本地运行需要安装 [`make`](https://www.gnu.org/software/make/)、[`go 1.19`](https://go.dev/doc/install)，之后编译并运行。



+ `make`: 安装依赖，并编译可在本机运行的 `Service`  与 `Client`  二进制文件（输出位置为: `./build/*`）。
  + 别名：`make build`。
  + `make amd64`: 与 `make`、`make build` 一致，区别是编译出的二进制文件仅能在 `amd64` 架构的机器上运行。
  + `make arm64`: 与 `make`、`make build` 一致，区别是编译出的二进制文件仅能在 `arm64` 架构的机器上运行。

+ `make build-service`: 安装依赖，并编译可在本机运行的 `Service`  二进制文件（输出位置为: `./build/service`）。
  + `make amd64-service`: 与 `make build-service` 一致，区别是编译出的二进制文件仅能在 `amd64` 架构的机器上运行。
  + `make arm64-service`: 与 `make build-service` 一致，区别是编译出的二进制文件仅能在 `arm64` 架构的机器上运行。

+ `make build-client`: 安装依赖，并编译可在本机运行的 `Client`  二进制文件（输出位置为: `./build/client`）
  + `make amd64-client`: 与 `make build-client` 一致，区别是编译出的二进制文件仅能在 `amd64` 架构的机器上运行。
  + `make arm64-client`: 与 `make build-client` 一致，区别是编译出的二进制文件仅能在 `arm64` 架构的机器上运行。

+ `make run-service`: 直接编译并运行 `Service` 。

+ `make run-client`: 直接编译并运行 `Client` 。
+ `make docker`: 编译可运行 `Client` 、`Service` 的 `docker`  镜像，`DOCKER_VERSION` 可指定标签。
  + 比如，需要 `1.0.0` 版本的，则可运行 `make docker DOCKER_VERSION=1.0.0` ，最后会编译出来 `cocytuselias2023/webproxy:service-1.0.0` 和 `cocytuselias2023/webproxy:client-1.0.0` 两个镜像。

+ `make docker-service`: 编译可运行 `Service` 的 `docker`  镜像，`DOCKER_VERSION` 可指定标签，`DOCKER_VERSION` 可指定标签。
  + 比如，需要 `1.0.0` 版本的，则可运行 `make docker-service DOCKER_VERSION=1.0.0` ，最后会编译出来 `cocytuselias2023/webproxy:service-1.0.0` 这个镜像。
+ `make docker-client`: 编译可运行 `Client` 的 `docker`  镜像，`DOCKER_VERSION` 可指定标签，`DOCKER_VERSION` 可指定标签。
  + 比如，需要 `1.0.0` 版本的，则可运行 `make docker-client DOCKER_VERSION=1.0.0` ，最后会编译出来 `cocytuselias2023/webproxy:client-1.0.0` 这个镜像。



### `Docker` 镜像

`Service` 和 `Client` 镜像名都是 `cocytuselias2023/webproxy`，镜像标签格式为 `type-tag`，`type` 分为 `Service` 和 `Client` 。

比如： `1.0.0` 版本的 `Service` 镜像就是 `cocytuselias2023/webproxy:service-1.0.0` 。`1.0.0` 版本的 `Client` 镜像就是 `cocytuselias2023/webproxy:client-1.0.0` 。



### 运行

先运行 `Service` 再运行 `Client`，运行 `Client` 前，需要在 `Client` 配置中指定 `Service` 服务地址。



## 配置

### 目录

本地配置目录在 `./config/`，`Service` 配置在 `./config/service.toml`，`Client` 配置在 `./config/client.toml`。

`Docker` 配置目录为 `/build/config`， `Service` 配置在 `/build/config/service.toml`，`Client` 配置在 `/build/config/client.toml`。



### Service 配置

`./example/service.toml` 为示例配置。

```toml
# 请求超时时间，单位秒。设置为 30，代表从收到请求开始算，30 秒内如果没收到 client 的响应，则直接向发送请求的客户端返回 408。
# 必须大于 0。
timeoutSecond = 30
# 监听服务地址。
addr = '0.0.0.0:8080'
# 最大请求数，当前正在处理中的请求最大数量。
# 必须大于 1000。
maxRequest = 3000
# 最大请求通道数，当前正在排队中的请求最大数量。
# 必须大于 100。
maxRequestChannel = 300
```



### Client 配置

`./example/client.toml` 为示例配置。

```toml
# service 的服务地址，必须以 ws:// 或 wss:// 为开头。
# 如想启用 wss:// ，需要在 service 前面添加 nginx 或 apisix，并配置 ssl 证书。
# url 路径必须是 /proxy/chan。
wsServiceUrl="ws://127.0.0.1:8080/proxy/chan"
# 请求超时时间，单位秒。设置为 15，代表从 client 向目标服务发起请求开始算，15 秒内如果没收到目标服务的响应，则终止请求。
timeoutSecond = 15
# 请求转发配置，这是第二个，匹配时也是按顺序匹配。如果这个匹配上了，就不执行后续的了。
[[routers]]
		# 路由匹配，以 api 为前缀的，使用正则表达式。
    path="^/api" 
		# 请求转发，会将这个服务地址的响应内容返回给 service。
    upstream="http://192.168.0.1:8080" 

# 请求转发配置，这是第二个，匹配时也是按顺序匹配。如果这个匹配上了，就不执行后续的了。
[[routers]]
		# 路由匹配，以 wap 为前缀的，使用正则表达式。
    path="^/wap"
    # 请求转发，会将这个服务地址的响应内容返回给 service。
    upstream="http://192.168.0.1:8081" 
    # 请求拷贝转发，这些服务地址只管请求，不管响应。
    copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]

# 请求转发配置，这是第三个，匹配时也是按顺序匹配。如果这个匹配上了，就不执行后续的了。
[[routers]]
		# 路由匹配，以 web 为前缀的，使用正则表达式。
    path="^/web"
    # 请求转发，会将这个服务地址的响应内容返回给 service。
    upstream="http://192.168.0.1:8081" 
    # 请求拷贝转发，这些服务地址只管请求，不管响应。
    copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]
    # 匹配到这个 route 时，需要按顺序执行的请求重写模块。
    # 这个是执行 stripPrefix 模块，参数为 4。既去除路径前四个字符，如果请求路径为 /web/user，那转发路径就只保留 /user。
    [[routers.rewriteModules]]
        name="stripPrefix"
        params=[4]

# 请求转发配置，这是第四个，匹配时也是按顺序匹配。如果这个匹配上了，就不执行后续的了。
[[routers]]
		# 路由匹配，以 oauth 为前缀的，使用正则表达式。
    path="^/oauth"
    # 请求转发，会将这个服务地址的响应内容返回给 service。
    upstream="http://192.168.0.11:8081"
    # 请求拷贝转发，这些服务地址只管请求，不管响应。
    copyStreams=["http://192.168.0.12:8082","http://192.168.0.13:8083"]
    # 匹配到这个 route 时，需要按顺序执行的请求重写模块。
    # 先使用 stripPrefix 重写, 再使用 stripSuffix 重写
    # 这个是先执行 stripPrefix 模块，参数为 6。再执行 stripSuffix 模块，参数为 4 。
    # 既先去除路径前 6 个字符，再去除路径最后 4 个字符。
    # 如果请求路径为 /oauth/user/login?a=1，那转发路径就只保留 /user/login。
    [[routers.rewriteModules]]
        name="stripPrefix"
        params=[6]
    [[routers.rewriteModules]]
        name="stripSuffix"
        params=[4]

# 请求转发配置，这是第五个，匹配时也是按顺序匹配。如果这个匹配上了，就不执行后续的了。
[[routers]]
		# 路由匹配，以 all 为前缀的，使用正则表达式。
    path="^/all"
    # 请求转发，会将这个服务地址的响应内容返回给 service。
    upstream="http://192.168.0.1:8081"
    # 请求拷贝转发，这些服务地址只管请求，不管响应。
    copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]
    # 匹配到这个 route 时，需要按顺序执行的请求重写模块。
    # 这个是执行 stripAll 模块，参数为 /hello/world。
    # 既，将整个路径替换为 /hello/world。
    # 如果请求路径为 /all/user/category，那直接重写为 /hello/world
    [[routers.rewriteModules]]
        name="stripAll"
        params=["/hello/world"]

# 请求转发配置，这是第五个，匹配时也是按顺序匹配。如果这个匹配上了，就不执行后续的了。
[[routers]]
		# 路由匹配，接受所有请求。
    path=".*" 
    # 请求转发，会将这个服务地址的响应内容返回给 service。
    upstream="http://192.168.0.1:8081"
    # 请求拷贝转发，这些服务地址只管请求，不管响应。
    copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]
    # 匹配到这个 route 时，需要按顺序执行的请求重写模块。
    # 这个是执行 stripAll 模块，参数为空字符串。
    # 既，将整个路径丢弃。
    # 如果请求路径为 /all/user/category，那直接重写为 /。
    [[routers.rewriteModules]]
        name="stripAll"
        params=[""]
```





## 模块

`Client` 支持编写请求重写模块。项目默认就只有三个模块，用来重写路由：

+ `stripPrefix`: 去除指定数量的前缀字符，参数只能为 `1` 个大于 `0` 的整数。
+ `stripSuffix`: 去除指定数量的后缀字符，参数只能为 `1` 个大于 `0` 的整数。
+ `stripAll`: 重写整个路径，参数只能为 `1` 个字符串，可以为空。



如果你想支持更复杂或者更高级的请求重写，你可以编写请求重写模块并注册。



### 存放目录

请求重写模块存放在 `./client/module` 下，每个目录是一个模块，每个模块中必须至少包含一个 `go` 文件，`go` 文件的 `package` 名称必须与所在的模块目录保持一致。



### 编写要求

一个完整的请求重写模块必须包含：

+ 一个名称为 `Module` 的 `struct`。
  + 内部的变量无要求，可以自定义属性用来存储信息
  + 需要有一个 `Handle(wsReq *constant.WsReq, params ...any) (wsReqRewrite *constant.WsReqRewrite, err error)` 方法，用来重写请求。这个方法的第一个参数是 `service` 发送给 `client` 的请求信息，这个方法的第二个参数是从匹配到的配置中拿到的参数。
  + 需要有一个 `Verify(params ...any) bool` 方法。主要是用来在 `client` 启动时，对相关的参数进行校验。
+ 一个 `Init` 方法，用来做模块初始化，并返回 `Module` 结构的指针。



### 编写步骤

1. 创建模块目录。
2. 创建模块 `go` 文件，`package` 名称需要与创建的模块目录相同。
3. 编写模块处理逻辑。
4. 执行 `make gen` 来注册模块。



### 示例

举个例子，如果我们想创建一个名为 `regex` 的模块，用来做正则路由重写，那我们需要：

1. 在 `./client/module/` 目录下创建一个名为 `regex` 的目录。现在我们有了 `./client/module/regex` 目录。

2. 在  `./client/module/regex` 目录下创建 `regex.go` 文件。

3. 文件内容参见下面的代码。

   ```go
   package regex
   
   import (
   	"webProxy/extern/constant"
   	"webProxy/extern/logger"
   )
   
   type Module struct {
   }
   
   // Init 初始化方法。
   func Init() (module *Module, err error) {
   	module = &Module{}
     
     // 编写你的初始化逻辑
     ... ...
   	
   	return module, nil
   }
   
   // Handle 处理方法。
   func (m *Module) Handle(wsReq *constant.WsReq, params ...any) (wsReqRewrite *constant.WsReqRewrite, err error) {
     	wsReqRewrite = &constant.WsReqRewrite{
   		Method: wsReq.Method,
   		Header: wsReq.Header,
   		Body:   wsReq.Body,
   		Path:   wsReq.Path,
   	}
     
     // 编写你的处理逻辑
     ... ...
     
   	return
   }
   
   // Verify 参数验证方法。
   func (m *Module) Verify(params ...any) bool {
     
     // 编写你的参数验证逻辑
     ... ...
   
   	return true
   }
   
   ```

4. 执行 `make gen` 来注册模块。

5. 之后，就可以在配置里使用下面的方式使用这个模块了。

   ```toml
   # 请求转发配置。
   [[routers]]
   		# 路由匹配，接受所有请求。
       path=".*" 
       # 请求转发，会将这个服务地址的响应内容返回给 service。
       upstream="http://192.168.0.1:8081"
       # 请求拷贝转发，这些服务地址只管请求，不管响应。
       copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]
       # 匹配到这个 route 时，需要按顺序执行的请求重写模块。
       # 这个是执行 stripAll 模块，参数为空字符串。
       # 既，将整个路径丢弃。
       # 如果请求路径为 /all/user/category，那直接重写为 /。
       [[routers.rewriteModules]]
           name="stripAll"
           params=[...] # 填写你的参数
   ```




## 维护者

[@eliassama](https://github.com/eliassama)