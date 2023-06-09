# WebProxy

![standard-readme compliant](https://static.elias.ink/resource/202306092335092.svg)
![standard-readme compliant](https://static.elias.ink/resource/202306092335497.svg)
![standard-readme compliant](https://static.elias.ink/resource/202306092335909.svg)
![standard-readme compliant](https://static.elias.ink/resource/202306092335041.svg)
![standard-readme compliant](https://static.elias.ink/resource/202306092335423.svg)



This is a lightweight intranet tunneling tool that only supports HTTP. It is not more powerful than [Ngrok](https://github.com/inconshreveable/ngrok) or [Frp](https://github.com/fatedier/frp). If you need a more powerful intranet tunneling and proxy tool, it is recommended to use other tools.

## Table of Contents
  - [What is WebProxy?](#What is WebProxy?)
  - [Running](#Running)
  - [Configuration](#Configuration)
  - [Modules](#Modules)
  - [Maintainer](#Maintainer)


## What is WebProxy?

WebProxy is a lightweight internal web proxy project that currently supports forwarding of the HTTP protocol.

It is not more powerful than [Ngrok](https://github.com/inconshreveable/ngrok) or [Frp](https://github.com/fatedier/frp). If you need a more powerful internal network tunneling proxy tool, it is recommended to use other tools.

Similar to general proxy services, WebProxy consists of two parts: the "Service" and the "Client". The "Service" is deployed on a public server (or any other server of your choice), while the "Client" is deployed on a machine that has simultaneous access to the internal services and the server where the "Service" is located.

Communication between the "Service" and the "Client" is done via WebSocket.


The overall request proxy process is as follows:

1. The "Service" starts an HTTP server to receive external network HTTP requests and acts as a WebSocket endpoint to accept WebSocket requests from the "Client".
2. The "Client" initiates a WebSocket request to the "Service" and establishes a WebSocket connection.
3. The "Service" receives incoming HTTP requests from the external network and forwards the request information to the "Client" through the WebSocket connection.
4. Upon receiving the request, the "Client" retrieves the specific forwarding service address and the rewritten request based on the configuration, and initiates an HTTP request.
5. After the "Client" sends the request, it forwards the response content to the "Service" through the WebSocket connection.
6. The "Service" uses the previous record to return the response content to the client that initiated the request.


## Running

### Local Make Compilation

To run it locally, you need to install [`make`](https://www.gnu.org/software/make/) and [`go 1.19`](https://go.dev/doc/install), and then compile and run it.

+ `make`: Install dependencies and compile the `Service` and `Client` binary files that can be run on your local machine (output location: `./build/*`).
  + Alias: `make build`.
  + `make amd64`: Same as `make` and `make build`, but the compiled binary files can only run on machines with the `amd64` architecture.
  + `make arm64`: Same as `make` and `make build`, but the compiled binary files can only run on machines with the `arm64` architecture.

+ `make build-service`: Install dependencies and compile the `Service` binary file that can be run on your local machine (output location: `./build/service`).
  + `make amd64-service`: Same as `make build-service`, but the compiled binary file can only run on machines with the `amd64` architecture.
  + `make arm64-service`: Same as `make build-service`, but the compiled binary file can only run on machines with the `arm64` architecture.

+ `make build-client`: Install dependencies and compile the `Client` binary file that can be run on your local machine (output location: `./build/client`).
  + `make amd64-client`: Same as `make build-client`, but the compiled binary file can only run on machines with the `amd64` architecture.
  + `make arm64-client`: Same as `make build-client`, but the compiled binary file can only run on machines with the `arm64` architecture.

+ `make run-service`: Compile and run the `Service` directly.

+ `make run-client`: Compile and run the `Client` directly.
+ `make docker`: Compile the `docker` images for running the `Client` and `Service`, with the `DOCKER_VERSION` specifying the tag.
  + For example, if you need version `1.0.0`, you can run `make docker DOCKER_VERSION=1.0.0`, and it will compile the `cocytuselias2023/webproxy:service-1.0.0` and `cocytuselias2023/webproxy:client-1.0.0` images.

+ `make docker-service`: Compile the `docker` image for running the `Service`, with the `DOCKER_VERSION` specifying the tag.
  + For example, if you need version `1.0.0`, you can run `make docker-service DOCKER_VERSION=1.0.0`, and it will compile the `cocytuselias2023/webproxy:service-1.0.0` image.
+ `make docker-client`: Compile the `docker` image for running the `Client`, with the `DOCKER_VERSION` specifying the tag.
  + For example, if you need version `1.0.0`, you can run `make docker-client DOCKER_VERSION=1.0.0`, and it will compile the `cocytuselias2023/webproxy:client-1.0.0` image.


### Docker Images

The Docker images for the "Service" and "Client" are both named `cocytuselias2023/webproxy`, and the image tags follow the format `type-tag`, where `type` can be either `Service` or `Client`.

For example, the image for the "Service" version `1.0.0` is `cocytuselias2023/webproxy:service-1.0.0`. The image for the "Client" version `1.0.0` is `cocytuselias2023/webproxy:client-1.0.0`.


### Run

To run the WebProxy, you need to first run the "Service" and then run the "Client". Before running the "Client", make sure to specify the address of the "Service" in the "Client" configuration.


## Configuration

Thank you for providing the configuration files for the WebProxy service. Here is a breakdown of the directory structure and the configuration files:

Directory structure:

- Local configuration directory: `./config/`
- Service configuration file: `./config/service.toml`
- Client configuration file: `./config/client.toml`
- Docker configuration directory: `/build/config`
- Service configuration file: `/build/config/service.toml`
- Client configuration file: `/build/config/client.toml`

Service Configuration (`./config/service.toml`):
```toml
# Request timeout in seconds. Set to 30, which means if no response is received from the client within 30 seconds from the request being received, a 408 error is returned to the client.
# Must be greater than 0.
timeoutSecond = 30
# Service listening address.
addr = '0.0.0.0:8080'
# Maximum number of requests, the maximum number of requests being processed concurrently.
# Must be greater than 1000.
maxRequest = 3000
# Maximum number of request channels, the maximum number of requests waiting in the queue.
# Must be greater than 100.
maxRequestChannel = 300
```

Client Configuration (`./config/client.toml`):
```toml
# Service URL of the WebProxy, must start with ws:// or wss://.
# If you want to use wss://, you need to add Nginx or Apache as a reverse proxy in front of the service and configure SSL certificates.
# The URL path must be /proxy/chan.
wsServiceUrl="ws://127.0.0.1:8080/proxy/chan"
# Request timeout in seconds. Set to 15, which means if no response is received from the target service within 15 seconds from the request being sent, the request is aborted.
timeoutSecond = 15

[[routers]]
    path="^/api"
    upstream="http://192.168.0.1:8080"

[[routers]]
    path="^/wap"
    upstream="http://192.168.0.1:8081"
    copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]

[[routers]]
    path="^/web"
    upstream="http://192.168.0.1:8081"
    copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]
    [[routers.rewriteModules]]
        name="stripPrefix"
        params=[4]

[[routers]]
    path="^/oauth"
    upstream="http://192.168.0.11:8081"
    copyStreams=["http://192.168.0.12:8082","http://192.168.0.13:8083"]
    [[routers.rewriteModules]]
        name="stripPrefix"
        params=[6]
    [[routers.rewriteModules]]
        name="stripSuffix"
        params=[4]

[[routers]]
    path="^/all"
    upstream="http://192.168.0.1:8081"
    copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]
    [[routers.rewriteModules]]
        name="stripAll"
        params=["/hello/world"]

[[routers]]
    path=".*"
    upstream="http://192.168.0.1:8081"
    copyStreams=["http://192.168.0.1:8082","http://192

.168.0.1:8083"]
    [[routers.rewriteModules]]
        name="stripAll"
        params=[""]
```

Please note that the configuration files you provided are examples (`./example/service.toml` and `./example/client.toml`). Make sure to modify the configurations according to your specific requirements before running the WebProxy service.


## Modules

`Client` supports writing custom request rewrite modules. By default, the project provides three modules for route rewriting:

- `stripPrefix`: Removes a specified number of prefix characters. The parameter must be a single integer greater than 0.
- `stripSuffix`: Removes a specified number of suffix characters. The parameter must be a single integer greater than 0.
- `stripAll`: Rewrites the entire path. The parameter must be a single string, which can be empty.

If you need more complex or advanced request rewriting capabilities, you can write and register your own request rewrite modules.

### Directory Structure

Request rewrite modules are stored in the `./client/module` directory. Each directory represents a module, and each module must contain at least one `go` file. The `package` name of the `go` file must match the name of the module's directory.

### Requirements

A complete request rewrite module must include:

- A `struct` named `Module`.
  - The internal variables can be customized to store information.
  - It should have a `Handle(wsReq *constant.WsReq, params ...any) (wsReqRewrite *constant.WsReqRewrite, err error)` method to handle request rewriting. The first parameter of this method is the request information sent from the `service` to the `client`, and the second parameter is the parameter obtained from the matched configuration.
  - It should have a `Verify(params ...any) bool` method, primarily used for validating parameters during `client` startup.
- An `Init` method to initialize the module and return a pointer to the `Module` structure.

### Writing Steps

1. Create a directory for your module.
2. Create a `go` file inside the module directory, ensuring that the `package` name matches the directory's name.
3. Write the module's processing logic.
4. Execute `make gen` to register the module.

### Example

Let's take an example of creating a module named `regex` for performing regex-based route rewriting. Here are the steps:

1. Create a directory named `regex` under `./client/module/`. Now we have the directory `./client/module/regex`.

2. Create a file named `regex.go` inside the `./client/module/regex` directory.

3. Put the following code in the `regex.go` file:

   ```go
   package regex
   
   import (
   	"webProxy/extern/constant"
   	"webProxy/extern/logger"
   )
   
   type Module struct {
   }
   
   // Init initializes the module.
   func Init() (module *Module, err error) {
   	module = &Module{}
     
     // Write your initialization logic here
     ... ...
   	
   	return module, nil
   }
   
   // Handle handles the request.
   func (m *Module) Handle(wsReq *constant.WsReq, params ...any) (wsReqRewrite *constant.WsReqRewrite, err error) {
     	wsReqRewrite = &constant.WsReqRewrite{
   		Method: wsReq.Method,
   		Header: wsReq.Header,
   		Body:   wsReq.Body,
   		Path:   wsReq.Path,
   	}
     
     // Write your processing logic here
     ... ...
     
   	return
   }
   
   // Verify performs parameter validation.
   func (m *Module) Verify(params ...any) bool {
     
     // Write your parameter validation logic here
     ... ...
   
   	return true
   }
   ```

4. Execute `make gen` to register the module.

5. Now you can use this module in the configuration using the following format:

   ```toml
   # Request forwarding configuration, this is the fifth one and will be matched in order.
   [[routers]]
   	# Route matching, accepts all requests.
       path=".*" 
       # Request forwarding, the response from this service address will be returned to the service.
       upstream="http://192.168.0.1:8081"
       # Request copy forwarding, these service addresses only handle requests, not responses.
       copyStreams=["http://192.168.0.1:8082","http://192.168.0.1:8083"]
       # Request rewrite modules to be executed in order when this route is matched.
       # In this example, the "stripAll" module is executed with an empty string parameter.
       # It will discard the entire path.
       # If the request path is "/all/user/category", it will be rewritten as "/".
       [[routers.rewriteModules]]
           name="stripAll"
           params=[...] # Fill in your parameters
   ```


## Maintainer

[@eliassama](https://github.com/eliassama)