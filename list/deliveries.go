package list

import (
	"github.com/gin-gonic/gin"
)

func Deliveries() *gin.Engine {
	//`gin.Default()` 是 Gin 框架中创建路由引擎的一种简便方法。它返回一个默认的 Gin 路由引擎实例
	//router := gin.Default()
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
