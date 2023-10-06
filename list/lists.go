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
