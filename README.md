tong - **桐** 是一个以 学习，研究 为目的的 Go Web 框架，遵循 GPL V3 协议。 

与 **桐** 框架配套的，是一本开源书籍《深入浅出 Go Web 编程》，在这本开源书籍中，我们尽可能全面，详细的介绍了与 Go Web 开发相关的知识。更重要的是，我们在这本书中完整的介绍了 **桐** 框架的设计思路，对源码进行了细致的解读。 

正如我们所说的， **桐** 是一个以 学习 为主要目的的 Web 框架，我们并不自负，也不盲目的想去创造一个全新的轮子。我们只是希望更多的人，不只是使用，而是去思考；我们只是希望的更多的人，不只是使用，而是去创造；我们希望通过对 **桐** 框架的使用，通过对《深入浅出 Go Web 编程》的阅读，能让更多的人发现，创造一个开源的框架，并不像想象中那么困难。很多之前刷题时研究过的算法，在工程中确有其实用性。 

当然，我们也希望，通过这个过程，你可以完善，细化自己 Web 开发的知识体系。 

我们也更希望，如果可能，你可以一起来完善 **桐** ，使它成为一个可以在生产环境中应用的完整框架。 

# A- 安装 

为了安装 tong package，首先你需要安装 Go 环境（version 1.11+ ），并设置好你的 Go workspace。 

随后，你可以使用下面的 Go 命令来安装 tong： 

```plain
go get github.com/ming3000/tong 
```
在你的代码中，通过下面的 import 语句来引入 tong： 
```plain
import "github.com/ming3000/tong" 
```
# A- 快速上手 

假定你的测试代码位于 demo.go 文件中。 

```go
$ cat ./demo.go 
package main 
import ( 
   "github.com/ming3000/tong" 
   "log" 
   "net/http" 
) 
func index(ctx *tong.Context) error { 
   return ctx.String(http.StatusOK, "hello world!") 
} 
func fun(ctx *tong.Context) error { 
   return ctx.String(http.StatusOK, "have fun!") 
} 
func main() { 
   t := tong.New() 
   t.GET("/", index) 
   t.GET("/fun", fun) 
   log.Fatalln(t.Start(":3000")) 
} 
```
运行 demo.go 文件： 
```plain
$ go run demo.go 
```
在浏览器中通过 localhost:3000/ 进行访问。 
# A- API 示例 

## a- GET & POST 

```go
func main() { 
   // create a tong instance 
   t := tong.New() 
    
   // register handler 
   t.GET("/someGet", somGetHandle) 
   t.POST("/somePost", somePostHandle) 
    
   // start the server on localhost:3000 
   log.Fatalln(t.Start(":3000")) 
} 
```
## a- Querystring parameters 

```plain
func main() { 
   // create a tong instance 
   t := tong.New() 
   // the handler math the url /user?name=python&age=30 
   t.GET("/user", func(c *tong.Context) error { 
      name := c.QueryString("name", "default name") 
      age := c.QueryInt("age", 3) 
       
      return c.String(http.StatusOK,  
                    fmt.Sprintf("hello, %s, %d", name, age)) 
   }) 
   log.Fatalln(t.Start(":3000")) 
} 
```
## a- Multipart/Urlencoded Form 

```plain
func main() { 
   // create a tong instance 
   t := tong.New() 
    
   t.POST("/user", func(c *tong.Context) error { 
      name := c.PostString("name", "default name") 
      age := c.PostInt("age", 3) 
      return c.String(http.StatusOK,  
                       fmt.Sprintf("hello,name:%s,age:%d", name, age)) 
   }) 
   log.Fatalln(t.Start(":3000")) 
} 
```

## a- String & JSON & ProtoBuf rendering 

返回 string 给客户端： 

```plain
String(code int, value string) error 
```
返回 json 结构体给客户端： 
```plain
Json(code int, value interface{}, indent string) error 
```
[todo] 
# A- Binding & Validate 

[todo] 

query string 

form 

json 

protobuf 

# A- 中间件 

[todo] 

日志 logger 

异常恢复 recover 

令牌桶限流 rating 

bloom-filter 过滤器 

跨域 

jwt 

# A- 定时任务 

[todo] 

# A- 日志 

[todo] 

# A- 缓存 

tong 的上下文 context 中，提供了 2 种缓存对象。它们都是并发安全的。 

|缓存对象名 |作用域 |实现策略 |
|:----|:----|:----|:----|
|RequestCache |该缓存的作用域为当前的 HTTP Request 的生命周期。 <br>在处理请求时，用作临时存储的对象，每次请求都会重设这个变量。 |LRU 链表 |
|GlobalCache |该缓存的作用域为整个 tong 应用程序的生命周期。 <br>该缓存的内容会被持久化到磁盘。下一次应用重启时，可以重复使用。 |LSM 文件 |

这两种缓存对象，都提供了下面的三种方法，供处理器程序使用： 

```plain
Set(key string, value interface{}) 
Get(key string) interface{} 
Del(key string) 
```
# A- 运行监控 

[todo] 

# 





