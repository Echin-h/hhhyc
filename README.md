# 项目结构
### 总体实现了快递的一些基本功能，发送删除更新列出查询的功能
## 载体形式
```go
type TODO struct {
	TrackingNumber string `json:"tracking_number"` // 跟踪号
	Time           string `json:"time"`            // 时间
	Location       string `json:"location"`        // 位置
	Recipient      string `json:"recipient"`       // 收件人
	Status         string `json:"status"`          // 状态
}

var todos []TODO // 存储所有的 TODO
```
用一个结构体切片来暂时储存数据
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
---
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
package list

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 更新TODO
func UpdateTodoByTrackingNumber(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")
	found := false

	for i, todo := range todos {
		if todo.TrackingNumber == trackingNumber {
			//设定一个bool值，可以更加明显的区分是否成功更新
			found = true
			var updatedTodo TODO
			//ShouldBindJson 和 BindJson  前者更高级，严谨

			//更新TODO
			if err := c.ShouldBindJSON(&updatedTodo); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}
			// 更新TODO
			//具体的意思就是 如果查到的单号数据不为空，则更新一下。
			if updatedTodo.Time != "" {
				todo.Time = updatedTodo.Time
			}
			if updatedTodo.Location != "" {
				todo.Location = updatedTodo.Location
			}
			if updatedTodo.Recipient != "" {
				todo.Recipient = updatedTodo.Recipient
			}
			if updatedTodo.Status != "" {
				todo.Status = updatedTodo.Status
			}
			todos[i] = todo
			SaveTodosToFile()
			break
		}
	}

	// http.StatusNotFound = 404
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "TODO not found"})
		return
	}
	//`gin.H{}` 是 map 类型， 写入键值对就行
	c.JSON(http.StatusOK, gin.H{"message": "TODO updated successfully"})
}
```
这个功能通过c.param实现传参，  使用一系列的if语句来判断如果结构体不是为空，则修改结构体，最后用savetodosfile 储存数据到文件
```go
URL /list/todo/{:tracking_number}
param输入  987654321
json输入  {       "tracking_number": "987654321",
        "time": "woshiren",
        "location": "nishiren",
        "recipient": "Jennifer Brown",
        "status": "Delivered"}
输出 ：{"message": "TODO updated successfully"} or {"error": "TODO not found"}
```
## 删除单号  delete.go
```go
package list
import (
	"github.com/gin-gonic/gin"
	"net/http"
)
// 删除TODO
func DeleteTodoByTrackingNumber(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")
	found := false

	for i, todo := range todos {
		if todo.TrackingNumber == trackingNumber {
			todos = append(todos[:i], todos[i+1:]...) // 从切片中删除匹配到的TODO
			found = true
			break
		}
	}

	if found {
		SaveTodosToFile()
		c.JSON(http.StatusOK, gin.H{"message": "TODO deleted successfully"})
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "TODO not found"})
		return
	}
}
```
具体功能的实现就是先查找，然后删除，然后保存到文件-- 细节就是可以通过一个found变量 来判断是否查找成功
```go
URL /list/todo/{:tracking_number}
c.param  tracking_number
响应 {"message": "TODO not found"} or {"message": "TODO deleted successfully"}
```
## 文件储存   file.go
```go
package list

import (
	"encoding/json"
	"log"
	"os"
)
// addTodo 添加 TODO
func AddTodo(todo *TODO) {
	todos = append(todos, *todo)
	SaveTodosToFile() // 保存到文件
}
// 创建文件
// 我通过重新创建文件 ， 并初始化文件， 把main 函数的长度缩减
func CreateFile() {
	//判断文件存不存在
	_, error := os.Stat("a.txt")
	// 文件不存在就创建一个  ， 存在的话就创建一个覆盖他
	file, err := os.OpenFile("a.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	stat, _ := file.Stat()  //检测状态
	if stat.Size() == 0 || error != nil {
		// 添加五个 TODO 项
		AddTodo(&TODO{
			TrackingNumber: "123456789",
			Time:           "2023-09-30 10:00",
			Location:       "Warehouse A",
			Recipient:      "John Doe",
			Status:         "In Transit",
		})
	..... 增加一些数字
	}
}

// loadTodosFromFile 从文件加载 TODO 列表
func LoadTodosFromFile() {
	file, err := os.Open("a.txt")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()      //defer关键字 ，延迟关闭
	decoder := json.NewDecoder(file) 
	err = decoder.Decode(&todos)    //将json数据解码 放到切片之中 ，通过todos的间接改变值
	if err != nil {
		log.Println("Error decoding todos:", err)
		return
	}
}

// saveTodosToFile 将 TODO 列表保存到文件
func SaveTodosToFile() {
	file, err := os.OpenFile("a.txt", os.O_RDWR|os.O_CREATE, 755)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)   //将切片中的数据编码为json格式，并通过encoder将其写入到文件中
	err = encoder.Encode(todos)
	if err != nil {
		log.Println("Error encoding todos:", err)  //日志
		return
	}
}
```
1. AddTodo 函数就是往切片中增加数据，同时储存到文件中
2. CreateFile 函数是创建文件， 如果没有文件或者文件为空，则会（创建文件）以及往文件中填充一些数据。
3. LoadTodosFromFile 函数 是从文件中加载数据，通过newdecoder 解码json 放入切片使用
4. saveTodosToFile 函数 是把现存的切片用newencoder编码成json 储存到文件中
---
---
# 创作过程
### 1-5 痛苦
一开始我是抱着高中的学习态度去学习的，作为一个 对go语言什么都没有涉及的新手，对边学边做的方法并不相信。 这几天每天都很痛苦，明明感觉学了很多东西，但是好像又对项目没啥作用
### 7-12 尝试
开始尝试边学边练（某个一起做这个项目的同学推荐），GPT，CSDN，Bilibli,想做什么搜什么，一开始完全是照抄别人的代码，问题是还实现不了，只能让GPT一遍又一遍地修改代码，
就这样我逐渐了解了代码的含义和怎样的思路去实现我想要的功能，也对Gin,restful,接口，APIfox有了更深的认识
### 13-15 完善
基本上有了方向， 我开始往代码中塞功能，听取学长的意见改代码，在完善的过程中，我注意到了许多代码中的细节（比如路程参数中要加冒号）  
同时我对APIfox的使用有了新的认识，至少不是单纯地点点发送按钮
### 收获
学习了restful风格， Gin框架， go基本语法 ，APIfox的简单使用 ， URL的了解， http协议
---
# 总结
1. 学的很辛苦 也很快乐。  幸苦的是没有思路的那几天，快乐的是代码跑出来的界面
2. 善用GPT  通过GPT我开始慢慢了解后端方面的简单代码实现
3. 边学边练，练什么学什么， 比只看课好用多了 ，印象还很深刻
4. 社团学长很友好， 群很活跃， 每次做不下去的时候 看看群活跃就感觉又可以坚持一下
5. 提高自主学习能力 ， 学长大多是给一个链接和文档让我们自己去学习，而不是直接手把手教导
6.提高问问题的问题，无论是问学长 ，还是问GPT
7. go语言真方便

**学的快乐最重要**




