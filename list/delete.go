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
