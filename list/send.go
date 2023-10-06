package list

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// addTodoFromForm 从表单数据添加 TODO
func AddTodoFromForm(c *gin.Context) {

	//另一种写法
	var entodo TODO
	if err := c.ShouldBindJSON(&entodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invailid request data"})
		return
	}
	AddTodo(&entodo)
	SaveTodosToFile()
	//http.StatusOK 和 http.StatusBadRequest  是两种状态码  也可以用200和400表示
	// 一个表示请求成功  ，一个表示请求存在错误
	c.JSON(http.StatusOK, gin.H{"message": "TODO added successfully"})
}

//一种写法
//`c.PostForm` 是 Gin 框架中的一个方法，用于获取客户端 POST 请求中的表单数据, 需要在body 中的post-data 输入，而不是json
/*trackingNumber := c.PostForm("tracking_number")
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
}*/
//其实这段代码也可以用ShouldBindJson实现，具体实现结果类似于update函数中的内容
