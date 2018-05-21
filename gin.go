package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"fmt"
)

func sampleDemo() {
	router := gin.Default()

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}

func parameterInPath() {
	router := gin.Default()
	// base
	// http://localhost:8080/user/mike
	router.GET("/user/:name", func(context *gin.Context) {
		name := context.Param("name")
		context.String(http.StatusOK, "Hello %s", name)
	})
	// collection path parameter test
	// http://localhost:8080/user/mike/hehe/haha?q=aa
	router.GET("/user/:name/*action", func(context *gin.Context) {
		name := context.Param("name")
		action := context.Param("action")
		message := name + " is " + action
		context.String(http.StatusOK, message)
	})

	// QueryString test
	// http://localhost:8080/welcome?name=mike&hobby=swim&hobby=climb
	router.GET("/welcome", func(context *gin.Context) {
		remoteAddr := context.Request.RemoteAddr
		requestURI := context.Request.RequestURI
		// query string
		firstname := context.DefaultQuery("firstname", "Guest")
		name := context.Query("name")
		queryArray := context.QueryArray("hobby")
		message := gin.H{
			"remoteAddr": remoteAddr,
			"requestURI": requestURI,
			"name": name,
			"hobby": queryArray,
			"firstname": firstname,

		}
		context.JSON(http.StatusOK, message)
	})
	router.Run()
}

//////////////////////////// FORM test//////////////////////
//message := context.GetPostForm()
//message := context.GetPostFormArray()
//message := context.DefaultPostForm()
//message := context.PostForm()
//message := context.PostFormArray()

//message := context.FormFile() // c.Request.FormFile(name)
//message := context.MultipartForm() // c.Request.ParseMultipartForm(c.engine.MaxMultipartMemory)


func multipartOrUrlencodeForm() {
	// form parameter test, urlencode and multipart is ok
	engine := gin.Default()
	engine.POST("/form_post", func(context *gin.Context) {
		defaultValue := context.DefaultPostForm("randomName", "defaultValue")
		name := context.PostForm("name")
		hobby := context.PostFormArray("hobby")
		context.JSON(http.StatusOK, gin.H{
			"defaultValue": defaultValue,
			"name": name,
			"hobby": hobby,
		})
	})
	engine.Run()
}

func singleFileUploadUseForm()  {
	engine := gin.Default()
	engine.POST("/form_post", func(context *gin.Context) {
		file, err := context.FormFile("file")
		if err == nil {
			log.Println(file.Filename)
			// Upload the file to specific dst.
			//dst := "/Users/mike/Desktop/" + file.Filename
			//context.SaveUploadedFile(file, dst)
			context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
		} else {
			context.String(http.StatusOK, fmt.Sprintf("upload error!!!"))
		}
	})
	engine.Run()
}

func multipleFilesUseForm() {
	// it also can handle one file
	engine := gin.Default()
	engine.POST("/post_form", func(context *gin.Context) {
		form, err := context.MultipartForm()
		if err != nil {
			context.String(http.StatusOK, "Upload error!!!!")
			return
		}
		files := form.File["upload[]"]
	})
	engine.Run()
}

//func contextTest() {
//context.XML()
			//context.Abort()
//
//}

func main() {
	//sampleDemo()
	//parameterInPath()
	multipartOrUrlencodeForm()
}