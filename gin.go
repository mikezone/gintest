package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"fmt"
	"os"
	"io"
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
			"name":       name,
			"hobby":      queryArray,
			"firstname":  firstname,
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
			"name":         name,
			"hobby":        hobby,
		})
	})
	engine.Run()
}

func singleFileUploadUseForm() {
	engine := gin.Default()
	engine.POST("/form_post", func(context *gin.Context) {
		file, err := context.FormFile("file")
		if err == nil {
			log.Println(file.Filename)
			// Upload the file to specific dst.
			dst := "/Users/mike/Desktop/" + file.Filename
			context.SaveUploadedFile(file, dst)
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
	engine.POST("/form_post", func(context *gin.Context) {
		form, err := context.MultipartForm()
		if err != nil {
			fmt.Println(err)
			context.String(http.StatusOK, "Upload error!!!!")
			return
		}
		files := form.File["upload[]"]
		for _, file := range files {
			log.Println(file.Filename)
			// Upload the file to specific dst.
			dst := "/Users/mike/Desktop/" + file.Filename + "000"
			context.SaveUploadedFile(file, dst)
		}
		context.String(http.StatusOK, fmt.Sprintf("%d files upload !!", len(files)))
	})
	engine.Run()
}

//func loginEndpoint(context *gin.Context) {
//	context.String(http.StatusOK, "login success!")
//}

func groupingRoutesTest() {
	engine := gin.Default()

	loginEndpoint := func(context *gin.Context) {
		context.String(http.StatusOK, "login success!")
	}

	v1 := engine.Group("/user")
	// add router for v1
	{
		v1.POST("/login", loginEndpoint)
	}
	// apply middleware on v1 group
	//v1.Use(AuthRequired())
	engine.Run()
}

func middlewareTest() {
	// Default With the Logger and Recovery middleware already attached
	//engine := gin.Default()

	// Blank Gin without middleware by default
	engine := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	engine.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	engine.Use(gin.Recovery())

	// Per route middleware, you can add as many as you desire.
	//engine.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := engine.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	AuthRequired := func(context *gin.Context) {
		name, ok := context.GetQuery("name")
		if !ok {
			context.JSON(http.StatusForbidden, gin.H{
				"authenticated": false,
			})
			context.Abort()
		} else {
			fmt.Println(name)
		}
	}
	authorized.Use(AuthRequired)
	{
		authorized.GET("/", func(context *gin.Context) {
			context.String(http.StatusOK, "ok")
		})
	}
	engine.Run()
}

func howToWriteLogFile()  {
	// Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

    // Use the following code if you need to write the logs to file and console at the same time.
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    engine := gin.Default()
    engine.GET("/ping", func(context *gin.Context) {
    	context.String(http.StatusOK, "pong")
	})
	engine.Run(":8080")
}

//func contextTest() {
//context.XML()
//context.Abort()
//}

func main() {
	//sampleDemo()
	//parameterInPath()
	//multipartOrUrlencodeForm()
	//singleFileUploadUseForm()
	//multipleFilesUseForm()
	//groupingRoutesTest()
	//middlewareTest()
	howToWriteLogFile()
}
