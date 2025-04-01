package main

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id int
	Title  string `json:"title" binding:"required"`
	Status bool `json:"status"`
}


func main() {
	todos := []Todo{
		{Id:1, Title: "Learn Go", Status: true},
		{Id:2, Title: "Build a project", Status: false},
		{Id:3, Title: "Learn Go language for tommorow", Status: true},
		{Id:4, Title: "Ata time for pes", Status: false},
	}
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.POST("/new_task", func(c *gin.Context){
		var data Todo
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		var id int
		if len(todos) > 0 {
			id = todos[len(todos)-1].Id + 1
		}else{
			id = 1
		}
		todos = append(todos, Todo{Id: id, Title: data.Title, Status: data.Status})


		c.JSON(http.StatusOK, gin.H{
			"todo": data,
		})
	})
	router.GET("/get_all", func(c *gin.Context){
		c.JSON(http.StatusOK, todos)
	}) 
	router.DELETE("/delete/:id", func(c *gin.Context){
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		newList := []Todo{}
		for _,todo := range todos {
			if todo.Id != int(id) {
				newList = append(newList, todo)
			}
		}
		todos = newList

		c.JSON(http.StatusOK, todos)
	})
	
	router.POST("/change_status/:id", func(c *gin.Context){
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"err":"Invalid id"})
			return
		}

		for i := range todos {
			if todos[i].Id == id {
				todos[i].Status = !todos[i].Status
			}
		}

		c.JSON(http.StatusOK, todos)
	})
	router.Run()
  }