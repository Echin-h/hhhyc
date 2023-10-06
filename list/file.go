package list

import (
	"encoding/json"
	"log"
	"os"
)

// TODO 结构体表示待办事项
type TODO struct {
	TrackingNumber string `json:"tracking_number"` // 跟踪号
	Time           string `json:"time"`            // 时间
	Location       string `json:"location"`        // 位置
	Recipient      string `json:"recipient"`       // 收件人
	Status         string `json:"status"`          // 状态
}

var todos []TODO // 存储所有的 TODO

// 下面是各种  文件的操作
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
	//通过这次条件语句 ，可以解决每次启动主函数都会有五组数据进入“a.txt”的问题
	//如果文件不存在，或者文件为空的时候，加进去这些数据
	stat, _ := file.Stat()
	if stat.Size() == 0 || error != nil {
		// 添加五个 TODO 项
		AddTodo(&TODO{
			TrackingNumber: "123456789",
			Time:           "2023-09-30 10:00",
			Location:       "Warehouse A",
			Recipient:      "John Doe",
			Status:         "In Transit",
		})
		AddTodo(&TODO{
			TrackingNumber: "987654321",
			Time:           "2023-09-30 14:00",
			Location:       "Warehouse B",
			Recipient:      "Jane Smith",
			Status:         "Delivered",
		})
		AddTodo(&TODO{
			TrackingNumber: "555555555",
			Time:           "2023-09-30 16:30",
			Location:       "Warehouse C",
			Recipient:      "Alice Johnson",
			Status:         "Out for Delivery",
		})
		AddTodo(&TODO{
			TrackingNumber: "111111111",
			Time:           "2023-09-30 12:30",
			Location:       "Warehouse D",
			Recipient:      "Bob Anderson",
			Status:         "In Transit",
		})
		AddTodo(&TODO{
			TrackingNumber: "999999999",
			Time:           "2023-09-30 18:00",
			Location:       "Warehouse E",
			Recipient:      "Emily Davis",
			Status:         "Pending",
		})
	}
}

// loadTodosFromFile 从文件加载 TODO 列表
func LoadTodosFromFile() {
	file, err := os.Open("a.txt")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)
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

	encoder := json.NewEncoder(file)
	err = encoder.Encode(todos)
	if err != nil {
		log.Println("Error encoding todos:", err)
		return
	}
}
