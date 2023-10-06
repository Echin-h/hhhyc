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
