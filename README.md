# 项目结构
### 总体实现了快递的一些基本功能，发送删除更新列出查询的功能
## 主函数 main.go
```go
package main
import (
	"todolist/list"
)
func main() {
	router := list.Deliveries()  //打开 该函数中的路由
	router.Run(":8080")  //监听请求
}
```
主函数导入了一个“todolist/list”的包（具体功能存储在这个包中）  
## 路由封装 -- 模块化   deliveries.go
```go
package list
import (
	"github.com/gin-gonic/gin"
)
func Deliveries() *gin.Engine {
	router := gin.Default()
	CreateFile()
	LoadTodosFromFile() // 从文件加载 TODO 列表
	delivery := router.Group("/list")
	{
		// /todo/:tracking_number 这边一定要加‘：’  不加冒号的话就无法实现连接
		delivery.GET("/todos", GetTodos)                                      // 获取所有 TODO
		delivery.GET("/todo/:tracking_number", GetTodoByTrackingNumber)       // 根据跟踪号查询 TODO
		delivery.POST("/todo", AddTodoFromForm)                               //发送单号TODO
		delivery.PUT("/todo/:tracking_number", UpdateTodoByTrackingNumber)    //更新单号TODO
		delivery.DELETE("/todo/:tracking_number", DeleteTodoByTrackingNumber) //删除单号TODO
	} //注意： 函数后面不用加（）
	return router
}
```
在deliveries 函数中 创建了一个Gin引擎 ， 并设置相应的路由 ，最后返回这个引擎给main.go    
个人感觉  这部分最难一块 是前期路径的理解和 一开始路径参数没加冒号导致APIfox 实现不了功能的悲伤

同时  对于GET POST PUT DELETE 这4项操作    同样是输入，也都可以通过json实现响应 ，为什么POST一定是增加（初学时百思不得其解）  
直到学习了 restful风格， 才知道这五个操作是一种规范，是让代码更具有有序性和规范性。
## 发送单号 send.go
```go
package list
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)
// addTodoFromForm 从表单数据添加 TODO
func AddTodoFromForm(c *gin.Context) {
	var entodo TODO
	if err := c.ShouldBindJSON(&entodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invailid request data"})
		return
	}
	AddTodo(&entodo)
	fmt.Println("ok")
	//http.StatusOK 和 http.StatusBadRequest  是两种状态码  也可以用200和400表示
	// 一个表示请求成功  ，一个表示请求存在错误
	c.JSON(http.StatusOK, gin.H{"message": "TODO added successfully"})
}
```
1. 首先这个函数他需要传入一个 gin.context类型的指针，包含了http请求的一些信息（设置响应状态，发送响应，请求参数）
2. shouldbindjson 是将json数据绑定到结构体上，说白了就是使用它 就可以 人为输入数据，同时bindjson和他差不多，但是没他严谨
3. c.json()是一个返回的函数，前面的表示状态，后面的表示返回的内容
总体来说， 这个功能是通过json渲染 ，我可以在body/json 输入数据  并通过AddTodo把数据加到文件之中。
```go
**实现案例**
{  
  "tracking_number": "ABC123",  
  "time": "2023-10-12 10:00:00",  
  "location": "New York",  
  "recipient": "John Doe",  
  "status": "Delivered"  
}
**URL** /list/todo
**返回**
{
    "message": "TODO added successfully"
}
```
```go
//一种写法
//`c.PostForm` 是 Gin 框架中的一个方法，用于获取客户端 POST 请求中的表单数据, 需要在body 中的post-data 输入，而不是json
trackingNumber := c.PostForm("tracking_number")
time := c.PostForm("time")
location := c.PostForm("location")
recipient := c.PostForm("recipient")
status := c.PostForm("status")
//这是 一个赋值的过程， 上面的代码得到time 的值 然后赋值给 todo.Time 为time
todo := TODO{
	TrackingNumber: trackingNumber,
	Time:           time,
	Location:       location,
	Recipient:      recipient,
	Status:         status,
}
```
这是一个可以通过form-data 传参的方式..  让我知道了form-data如何传参。
## 列出单号 lists.go
```go
package list
import (
	"github.com/gin-gonic/gin"
	"net/http"
)
// getTodos 获取所有 TODO
func GetTodos(c *gin.Context) {
	LoadTodosFromFile() // 重新加载文件数据
	c.JSON(http.StatusOK, todos)
}
```
这个功能就是从文件中加载数据，（数据是结构体切片）然后用c.json返回即可
## 查询单号 check.go
```go
package list

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// getTodoByTrackingNumber 根据跟踪号查询 TODO
func GetTodoByTrackingNumber(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")
	for _, todo := range todos {
		if todo.TrackingNumber == trackingNumber {
			c.JSON(http.StatusOK, todo)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "TODO not found"})
}
```
这个功能的实现 让我学会了路径参数的输入，以前我是直接在URL上面改写路径参数，也没有管params中query和paht参数的选项有什么用
同时功能的实现 就是从文件中现有切片（数据是储存在切片中）中的数据与我输入的数据相匹配，如果匹配则输出匹配的那个结构体，否则，输出error
```go
URL /list/todo/:{tracking_number}
输入 123456789
输出{
    "tracking_number": "987654321",
    "time": "",
    "location": "",
    "recipient": "Jane Smith",
    "status": "Delivered"
}
输入 55664646
输出{
    "error": "TODO not found"
}
```
## 更新单号  renewal.go
```go

```
分别设有 发送单号（send.go）  删除单号（delete.go）  更新单号（renewal.go）  列出单号（lists.go） 查询单号（check.go） 文件操作（file.go） 路由封装（deliveries.go） 六个部分
所有部分 集结与 list 包中 ， 在 main.go 中可以通过import("todolist/list")来调用  
send.go /list/todo          body/json 输入
delete.go /list/todo/{tracking_number}    params paths参数输入
renewal.go /list/todo/{tracking_number}   params paths参数输入  body/json 输入
check.go  /list/todo/{tracking_number}    params paths 参数输入
lists.go /list/todo    点击发送即可

/二创作过程
作为一个 对go语言什么都没有涉及的新手， 面对demo中的gin框架，http协议，go语法.. 陌生的知识 ， 虽然学的很慢， 但是一直在学

首先我觉得自己前期还是没有理解边做边练的意思
过分的追求 “全面掌握go语法的知识”  ，想要把所有的知识都学过以后才开始动手
到了第五天 ，我发现我除了一个demo 啥都没有
同时我对 Gin json apifox 的了解 甚少

到了第六天或者第七天的时候
我选择了让 GPT来帮我写项目 ，我来负责看项目， 理解项目
前几天是囫囵吞枣度过的，  一知半解的记住了形式  要json传递  要使用bindjson（&todo）之类的/某个的代码形式 差不多能实现某个内容
渐渐地到了十多天的时候， 写 删除todo更新todo的时候，我已经能够自己思路较为清晰地 打出自己想要打出的代码（虽然后期还是得靠GPT）
同时  前五六天 疯狂（没看进去脑子）的知识点 模模糊糊地加深了印象。
我也能够 理解 apifox 中 params json post-data 的意思

最后几天 我给代码进行了模块化 ，对学长提出的改进建议进行了改进。

说实话 在第六天第七天，那个时候 我没有一刻不是想放弃这个面试任务的---我觉得没有方向，觉得我虽然学了知识，但是没有用处
但是 我 一直想着方丈在群里说的话： 只要你肯学，杭电助手就欢迎。
我觉得  这句话一直激励着我 哪怕坐在电脑前没有方向的时候，哪怕每天晚上学到很晚也不知道干啥的时候，我还是会去寻找 这个任务的方向， 我需要提升的方向，
同时
我一点也不认可 学的快乐最重要
但我认可 “爱飞的鸟”  学长说的   敲代码敲得快乐最重要
我自己也尝到了-- 代码运行起来后的快乐最重要
/三 感谢
对于指导者（我一直是对此保持感激心态的）
1.GPT  和他对话中， 极大地削减了我心中的苦闷， 还能够快速地获得新知识 。是他在前期我啥都不知道的时候 无声给予他所知道的代码。 可以说 每天 和 GPT深情交流的时间是最多的  
由此也提高我问问题的水平  ，非常感谢。
2.CSDN  茫茫文章中寻找对自己有意义的文章 ，在CSDN中可以获取 自己想要明白的一部分知识 ，（但终究有些文字过于生涩专业 难以理解 惭愧）非常感谢
3.B站   无论是老郭的go语言 还是 不知名UP 的GIN框架搭建介绍 ， json 路由  等 视频的讲解，即使不能一点就通 ，还是能够让我或多或少对此有点了解 ， 非常感谢
4.学长   我觉得 学长永远都是 激励着让我不要放弃的源泉。  做任务的时候 （哪怕再痛苦的时候）  依旧可以从学长中感受到他们的温暖。 
言传身教   从学长身上  我可以明显的感受到他们 想要培养我们自主学习操作的能力 -- （一开始我还抱怨学长不直接告诉答案  非常抱歉）  
一起做任务的人有很多 ，学长自身的事情也有很多 ，但他们仍旧能够抽出时间来 帮助我 ，非常感谢。  感谢方丈，感谢爱飞的鸟，感谢iyear
5.任务简介  说实话刚刚开始拿到任务的时候 我是很蒙蔽的 ，都不知道要干啥 ，但是任务简介中 向我们推送了许多权威的视频和文档 ，可以让我们不用再去碰壁同时更快地理解后端的意义和知识点， 非常感谢
6.家人   有姐姐和哥哥每天听我抱怨，你们总是能够在我想要放弃的时候 疏散我心中的怨气， 让我整理好心态继续向前，非常感谢
7. 自己  感谢自己 能够活着做到这个项目，虽然项目只完成了 一些小功能，没有数据库也没有接口鉴权 ，但我觉得我还是迈出了一步 ，至少我没有放弃，一直在努力学习。


/四 收获
首先我了解基本的语法的使用
对gin框架的基本 认识
restful 风格 --一种很简约有规则的风格
对URL 接口（一种约定） http协议的了解
对APIfox params query json 的使用
对路径的{：}的了解
以及一些go语言库函数的使用 （我觉得这方面的认识真的得靠实践 而不是做笔记）





