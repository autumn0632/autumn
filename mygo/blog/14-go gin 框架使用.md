# 一、 gin框架简介

1. [Gin](https://github.com/gin-gonic/gin) 是用 Go 编写的一个 Web 应用框架，他有更好的性能和更快的路由，方便灵活的中间件，强大的gin.Context

2. 初始化过程

   > - 创建一个 Engine 对象
   > - 注册中间件
   > - 注册路由（组）

3. 快速运行：

   ```go
   package main
   
   import (
       "github.com/gin-gonic/gin"
       "net/http"
   )
   
   func main() {
       // 初始化引擎
       engine := gin.Default()
       // 注册一个路由和处理函数
     	engine.Any("/", func(c *gin.Context){
           context.String(http.StatusOK, "hello, world")
     	})
       // 绑定端口，然后启动应用
       engine.Run()
   }
   
   /**
   * 根请求处理函数
   * 所有本次请求相关的方法都在 context 中，完美
   * 输出响应 hello, world
   */
   func WebRoot(context *gin.Context) {
       context.String(http.StatusOK, "hello, world")
   }
   ```

4. 支持动态路由

   > 需求：例如 `/user/:id`，通过调用的 url 来传入不同的 id .

5. 路由组支持

   ```go
   // 省略的代码 ...
   
   func main() {
       router := gin.Default()
   
       // 定义一个组前缀
     	// /v1/login 就会匹配到这个组
       v1 := router.Group("/v1")
       {
           v1.POST("/login", loginEndpoint)
           v1.POST("/submit", submitEndpoint)
           v1.POST("/read", readEndpoint)
       }
   
       // 定义一个组前缀
     	// 不用花括号包起来也是可以的。上面那种只是看起来会统一一点。看你个人喜好
       v2 := router.Group("/v2")
       v2.POST("/login", loginEndpoint)
       v2.POST("/submit", submitEndpoint)
       v2.POST("/read", readEndpoint)
   
       router.Run(":8080")
   }
   ```

6. 中间件（Middleware）

   > 1. 中间件的写法和路由的 Handler 几乎是一样的，只是多调用 `c.Next()`。
   > 2. c.Next()可以用来流程控制

7. 参数

   > 1. url参数
   >
   >    **c.Query()、c.DefaultQuery()**
   >
   >    ```go 
   >    // 注册路由和Handler
   >        // url为 /welcome?firstname=Jane&lastname=Doe
   >        router.GET("/welcome", func(c *gin.Context) {
   >            // 获取参数内容
   >            // 获取的所有参数内容的类型都是 string
   >            // 如果不存在，使用第二个当做默认内容
   >            firstname := c.DefaultQuery("firstname", "Guest")
   >            // 获取参数内容，没有则返回空字符串
   >            lastname := c.Query("lastname") 
   >    
   >            c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
   >        })
   >    ```
   >
   > 2. 动态路由参数：
   >
   >    **c.Param()**
   >
   >    ```go
   >    // Param returns the value of the URL param.
   >    // It is a shortcut for c.Params.ByName(key)
   >    	router.GET("/user/:id", func(c *gin.Context) {
   >             // a GET request to /user/john
   >             id := c.Param("id") // id == "john"
   >         })
   >    ```
   >
   >    
   >
   > 3. 表单和Body参数（Multipart/Urlencoded Form）
   >
   >    **c.PostFrom()**
   >
   >    ```go
   >    router.POST("/form_post", func(c *gin.Context) {
   >            // 获取post过来的message内容
   >            // 获取的所有参数内容的类型都是 string
   >            message := c.PostForm("message")
   >            // 如果不存在，使用第二个当做默认内容
   >            nick := c.DefaultPostForm("nick", "anonymous")
   >    
   >            c.JSON(200, gin.H{
   >                "status":  "posted",
   >                "message": message,
   >                "nick":    nick,
   >            })
   >        })
   >    ```
   >
   >    

8. 数据绑定

   > 使用 `c.ShouldBindQuery`方法，可以自动绑定 Url 查询参数到 `struct`.

9. 输出响应

   > 1. string:
   >
   >    ```go
   >    func Handler(c *gin.Context) {
   >        // 使用 String 方法即可
   >        c.String(200, "Success")
   >    }
   >    ```
   >
   > 2. json、xml、yaml
   >
   >    使用`gin.H{}`
   >
   >     *// 会输出头格式为 application/json; charset=UTF-8 的 json 字符串* 
   >
   >     c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK}) 
   >

10. 请求声明周期

   golang原生为web而生而提供了完善的功能，gin能做的事情也是把 `ServeHTTP(ResponseWriter, *Request)` 做得高效、友好。一个请求来到服务器了，`ServeHTTP` 会被调用，gin做的事情包括：

   > * 路由，找到handle
   > * 将请求和响应用Context包装起来供业务代码使用
   > * 依次调用中间件和处理函数
   > * 输出结果

# 二、gin高级

## 1. Context 结构体

`context.go`  包含主体功能代码，包含了** HTTP 请求上下文全部处理过程, `request` 和 `response` 两部分.**

**Context作为一个数据结构在中间件中传递本次请求的各种数据、管理流程，进行响应**

```go
// Context作为一个数据结构在中间件中传递本次请求的各种数据、管理流程，进行响应
// context.go:40
type Context struct {
    // ServeHTTP的第二个参数: request
    Request   *http.Request

    // 用来响应 
    Writer    ResponseWriter
    writermem responseWriter

    // URL里面的参数，比如：/xx/:id  
    Params   Params
    // 参与的处理者（中间件 + 请求处理者列表）
    handlers HandlersChain
    // 当前处理到的handler的下标
    index    int8

    // Engine单例
    engine *Engine

    // 在context可以设置的值
    Keys map[string]interface{}

    // 一系列的错误
    Errors errorMsgs

    // Accepted defines a list of manually accepted formats for content negotiation.
    Accepted []string
}

// response_writer.go:20
type ResponseWriter interface {
    http.ResponseWriter //嵌入接口
    http.Hijacker       //嵌入接口
    http.Flusher        //嵌入接口
    http.CloseNotifier  //嵌入接口

    // 返回当前请求的 response status code
    Status() int

    // 返回写入 http body的字节数
    Size() int

    // 写string
    WriteString(string) (int, error)

    //是否写出
    Written() bool

    // 强制写htp header (状态码 + headers)
    WriteHeaderNow()
}

// response_writer.go:40
// 实现 ResponseWriter 接口
type responseWriter struct {
    http.ResponseWriter
    size   int
    status int
}


type errorMsgs []*Error


// 每当一个请求来到服务器，都会从对象池中拿到一个context。其函数有：
// **** 创建
reset()                 //从对象池中拿出来后需要初始化
Copy() *Context         //克隆，用于goroute中
HandlerName() string    //得到最后那个处理者的名字
Handler()               //得到最后那个Handler

// **** 流程控制
Next()                  // 只能在中间件中使用，依次调用各个处理者
IsAborted() bool    
Abort()                 // 废弃
AbortWithStatusJson(code int, jsonObj interface{})
AbortWithError(code int, err error) *Error

// **** 错误管理
Error(err error) *Error // 给本次请求添加个错误。将错误收集然后用中间件统一处理（打日志|入库）是一个比较好的方案

// **** 元数据管理 ****
Set(key string, value interface{})  //本次请求用户设置各种数据 (Keys 字段)
Get(key string)(value interface{}, existed bool)
MustGet(key string)(value interface{})
GetString(key string) string
GetBool(key string) bool
GetInt(key string) int
GetInt64(key string) int64
GetFloat64(key string) float64
GetTime(key string) time.Time
GetDuration(key string) time.Duration
GetStringSlice(key string) []string
GetStringMap(key string) map[string]interface{}
GetStringMapString(key string) map[string]string
GetStringMapStringSlice(key string) map[string][]string

// **** 输入数据 ****
//从URL中拿值，比如 /user/:id => /user/john
Param(key string) string    

//从GET参数中拿值，比如 /path?id=john
GetQueryArray(key string) ([]string, bool)  
GetQuery(key string)(string, bool)
Query(key string) string
DefaultQuery(key, defaultValue string) string
GetQueryArray(key string) ([]string, bool)
QueryArray(key string) []string

//从POST中拿数据
GetPostFormArray(key string) ([]string, bool)
PostFormArray(key string) []string 
GetPostForm(key string) (string, bool)
PostForm(key string) string
DefaultPostForm(key, defaultValue string) string

// 文件
FormFile(name string) (*multipart.FileHeader, error)
MultipartForm() (*multipart.Form, error)
SaveUploadedFile(file *multipart.FileHeader, dst string) error

// 数据绑定
Bind(obj interface{}) error //根据Content-Type绑定数据
BindJSON(obj interface{}) error
BindQuery(obj interface{}) error

//--- Should ok, else return error
ShouldBindJSON(obj interface{}) error 
ShouldBind(obj interface{}) error
ShouldBindJSON(obj interface{}) error
ShouldBindQuery(obj interface{}) error

//--- Must ok, else SetError
MustBindJSON(obj interface{}) error 

ClientIP() string
ContentType() string
IsWebsocket() bool

// **** 设置输出数据 ****
Status(code int)            // 设置response code
Header(key, value string)   // 设置header
GetHeader(key string) string

GetRawData() ([]byte, error)

Cookie(name string) (string, error)     // 设置cookie
SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)

Render(code int, r render.Render)      // 数据渲染
HTML(code int, name string, obj interface{})    //HTML
JSON(code int, obj interface{})                 //JSON
IndentedJSON(code int, obj interface{})
SecureJSON(code int, obj interface{})
JSONP(code int, obj interface{})                //jsonp
XML(code int, obj interface{})                  //XML
YAML(code int, obj interface{})                 //YAML
String(code int, format string, values ...interface{})  //string
Redirect(code int, location string)             // 重定向
Data(code int, contentType string, data []byte) // []byte
File(filepath string)                           // file
SSEvent(name string, message interface{})       // Server-Sent Event
Stream(step func(w io.Writer) bool)             // stream


// **** 实现  context.Context 接口(GOROOT中)
```



## 2. 路由（Router）

gin 的高效，很大一部分是说其路由效率。

```go 
// 路由API
// routergroup.go:20
type IRoutes interface {
    Use(handlers ...HandlerFunc) IRoutes

    Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes
    Any(relativePath string, handlers ...HandlerFunc) IRoutes
    GET(relativePath string, handlers ...HandlerFunc) IRoutes
    POST(relativePath string, handlers ...HandlerFunc) IRoutes
    DELETE(relativePath string, handlers ...HandlerFunc) IRoutes
    PATCH(relativePath string, handlers ...HandlerFunc) IRoutes
    PUT(relativePath string, handlers ...HandlerFunc) IRoutes
    OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes
    HEAD(relativePath string, handlers ...HandlerFunc) IRoutes

    StaticFile(relativePath, filepath string) IRoutes
    Static(relativePath, root string) IRoutes
    StaticFS(relativePath string, fs http.FileSystem) IRoutes
}

// routergroup.go:15
type IRouter interface {
    IRoutes
    Group(string, ...HandlerFunc) *RouterGroup
}
```

