package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Item struct {
	ID    int `json:"id"`
	Name  string `json:"name"`
	Created_at time.Time `json:"date"`
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.MaxMultipartMemory = 8 << 20  
	router.POST("/upload", func(c *gin.Context) {
		jsonFileName := "data.json"	
		var items []Item
		if _,err := os.Stat(jsonFileName); os.IsNotExist(err){
			createFile, err := os.Create(jsonFileName)
			if err != nil {
				log.Fatal(err)
				return
			}
			defer createFile.Close()
		} else if err != nil {
			log.Println("Error checking file:", err)
		} else {
			data,err := os.ReadFile(jsonFileName)
			if err != nil {
				log.Printf("Failed to read file: %v\n", err)
			}
			err = json.Unmarshal(data, &items)
			if err != nil {
				log.Printf("Failed to parse JSON: %v\n", err)
			}
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to get file")
			return
		}

		var newID int
		if len(items) == 0 {
			newID = 1
		} else {
			newID = items[len(items)-1].ID + 1
		}
		addItem := Item{ID: newID, Name: file.Filename, Created_at: time.Now()}
		items = append(items, addItem)

		localFile, err := os.OpenFile(jsonFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("Failed to read file: %v\n", err)
			return
		}
		defer localFile.Close()

		jsonItems, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
		log.Printf("Failed to marshal JSON: %v\n", err)
		}

		_, err = localFile.WriteString(string(jsonItems))
		if err != nil {
		log.Printf("Failed to write data: %v\n", err)
		}


		
		c.SaveUploadedFile(file, "./files/" + file.Filename)
		c.JSON(http.StatusOK, gin.H{
			"message": "uploaded",
		})
	})
	router.GET("/get_uploaded_list", func(c *gin.Context){
		jsonFileName := "data.json"	

		localFile, err := os.ReadFile(jsonFileName)
		if err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"message":"something went wrong in finding file"})
			return
		}
		c.Data(http.StatusOK, "application/json", localFile)
	})
	router.GET("/download/:filename", func(c *gin.Context){
		filename := c.Param("filename")
		filePath := "./files/" + filename

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		c.File(filePath)
	})
	router.Run()
}
