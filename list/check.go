package list

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// getTodoByTrackingNumber 根据跟踪号查询 TODO
func GetTodoByTrackingNumber(c *gin.Context) {
	//`c.Param` 是 Gin 框架中的一个方法， 就是添加参数 APIFox中的params
	//`c.query` 也是， 也是添加参数 ，APIFox 中在params 上面
	//在 Gin 框架中，可以在路由路径中使用冒号（:）来定义参数占位符。
	//总而言之就是传入具体的数值（在APIfox中） 然后会显现出特定的内容
	trackingNumber := c.Param("tracking_number")

	for _, todo := range todos {
		if todo.TrackingNumber == trackingNumber {
			c.JSON(http.StatusOK, todo)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "TODO not found"})
}
