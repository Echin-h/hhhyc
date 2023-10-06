package main

import (
	"todolist/list"
)

func main() {

	//一些 增删改查的todo 小功能
	router := list.Deliveries()
	router.Run(":8080")
}
