package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"log"
)

func httpMethod() {
	router := gin.Default()
	router.GET("/get/:name", func(c *gin.Context) {
		//c.Query()
		//c.Param()
		//c.FormFile()
		name := c.Param("name")
		c.String(http.StatusOK, "hello %s", name)
	})
	router.POST("/post", func(c *gin.Context) {
		//c.GetPostForm()
		//c.PostForm()
		name, ok := c.GetPostForm("name")
		c.String(http.StatusOK, "%s, %v", name, ok)
	})
	router.PUT("/put", func(c *gin.Context) {
		multipartForm, err := c.MultipartForm()
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusBadRequest, "error")
		} else {
			fmt.Println(multipartForm.Value["name"][0])
			// Multipart form
			files := multipartForm.File["upload[]"]

			for _, file := range files {
				log.Println(file.Filename)
				// Upload the file to specific dst.
				// c.SaveUploadedFile(file, dst)
			}
			c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
			////////////////////////////////////////
			// Formfile
			// Upload the file to specific dst.
			file, _ := c.FormFile("file")
			dst := "/Users/mike/Desktop/" + file.Filename
			c.SaveUploadedFile(file, dst)
			log.Println(file.Filename)
			c.String(http.StatusOK, fmt.Sprintf("%s", file.Filename))
		}
	})
	router.GET("/benchmark", func(c *gin.Context) {
	})
	router.Run()
}

func sampleDemo() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func main() {
	//sampleDemo()
	httpMethod()
}
