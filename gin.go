package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"fmt"
	"os"
	"io"
	"time"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"github.com/gin-gonic/gin/binding"
	"html/template"
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

func howToWriteLogFile() {
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

// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func bookableDate(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string, ) bool {
	if date, ok := field.Interface().(time.Time); ok {
		today := time.Now()
		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
			return false
		}
	}
	return true
}

type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

type myForm struct {
	Colors []string `form:"colors[]"`
}

func modelBindingTest() {

	engine := gin.Default()

	// post body is json-formated string
	/**
	** both of below are ok

	Content-Type: text/plain; charset=utf-8
	{"user": "Mike", "password": "123"}

	Content-Type: application/json; charset=utf-8
	{"user":"Mike","password":"123"}
	*/
	engine.POST("/loginJSON", func(context *gin.Context) {
		var json Login
		if err := context.ShouldBindJSON(&json); err == nil {
			if json.User == "Mike" && json.Password == "123" {
				context.JSON(http.StatusOK, gin.H{"message": "login success!!"})
			} else {
				context.JSON(http.StatusOK, gin.H{"message": "unauthorized"})
			}
		} else {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}
	})

	// post body is HTML form
	/*
	both of below are ok

	Content-Type: multipart/form-data; charset=utf-8; boundary=__X_PAW_BOUNDARY__

	Content-Type: application/x-www-form-urlencoded; charset=utf-8
	*/
	engine.POST("/loginForm", func(context *gin.Context) {
		var form Login
		if err := context.ShouldBind(&form); err == nil {
			if form.User == "Mike" && form.Password == "123" {
				context.JSON(http.StatusOK, gin.H{"message": "ok"})
			} else {
				context.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			}
		} else {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	getBookable := func(context *gin.Context) {
		var b Booking
		if err := context.ShouldBindWith(&b, binding.Query); err == nil {
			context.JSON(http.StatusOK, gin.H{"message": "book is valid"})
		} else {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	// test: http "localhost:8080/bookable?check_in=2018-05-24&check_out=2018-10-02"
	engine.GET("/bookable", getBookable)

	// only bind query string
	startPage := func(context *gin.Context) {
		var person Person
		if context.ShouldBindQuery(&person) == nil {
			log.Println("====== Only Bind By Query String =====")
			log.Println(person.Name)
			log.Println(person.Address)
		}
		context.String(http.StatusOK, "success")
	}
	// for test: http 'localhost:8080/testing?name=Mike&address=hehe'
	// only bind querystring, both of below have no effect
	/*
	application/x-www-form-urlencoded:
	name=Mike&address=123

	multipart/form-data
	--__X_PAW_BOUNDARY__
	Content-Disposition: form-data; name="name"

	Mike
	--__X_PAW_BOUNDARY__
	Content-Disposition: form-data; name="address"

	123
	--__X_PAW_BOUNDARY__--
	*/
	engine.Any("/testing", startPage)

	// bind query or post data : ShoudBind
	bindQueryOrPostData := func(context *gin.Context) {
		var person Person
		if context.ShouldBind(&person) == nil {
			log.Println(person.Name)
			log.Println(person.Address)
			log.Println(person.Birthday)
		}
		context.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
	// for test:
	/*
	QueryString:
	http 'localhost:8080/bindQueryOrForm?name=Mike&address=hehe&birthday=2018-05-20'

	Form:
	application/x-www-form-urlencoded
	multipart/form-data
	*/
	engine.Any("/bindQueryOrForm", bindQueryOrPostData)

	// bind checkbox
	handlerCheckBox := func(context *gin.Context) {
		var form myForm
		context.ShouldBind(&form)
		context.JSON(http.StatusOK, gin.H{
			"colors": form.Colors,
		})
	}
	// http 'localhost:8080/bindCheckBox?colors[]=red&colors[]=green'
	engine.GET("/bindCheckBox", handlerCheckBox)
	// test:  post checkbox
	// application/x-www-form-urlencoded
	// http -f POST localhost:8080/bindCheckBox 'colors[]=red' 'colors[]=green'
	// multipart/form-data
	// http -f POST localhost:8080/bindCheckBox 'colors[]=red' 'colors[]=green' 'colors[]=cyan' file@~/Desktop/a.txt
	engine.POST("/bindCheckBox", handlerCheckBox)

	// run
	engine.Run()
}

//XML, JSON and YAML rendering
func responseContentTypeTest() {
	engine := gin.Default()

	// gin.H is a shortcut for map[string]interface{}
	// test: http localhost:8080/someJSON
	engine.GET("/someJSON", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	// test: http localhost:8080/moreJSON
	engine.GET("/moreJSON", func(context *gin.Context) {
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// Note that msg.Name becomes "user" in the JSON
		context.JSON(http.StatusOK, msg)
	})

	// test: http localhost:8080/someXML
	engine.GET("/someXML", func(context *gin.Context) {
		context.XML(http.StatusOK, gin.H{"message": "ok"})
	})
	// test: http localhost:8080/someYAML
	engine.GET("/someYAML", func(context *gin.Context) {
		context.YAML(http.StatusOK, gin.H{"message": "ok"})
	})
	engine.Run()
}

func serveStaticFiles() {
	engine := gin.Default()
	engine.Static("/static", "./static")                         // http http://localhost:8080/static/1.html
	engine.StaticFS("/statictest", http.Dir("./fakestatic"))     // http http://localhost:8080/statictest/1.html
	engine.StaticFile("/favicon.ico", "./resources/favicon.ico") // http://localhost:8080/favicon.ico
	engine.Run()
}

func servingDataFromReader() {
	engine := gin.Default()
	engine.GET("/someDataFromReader", func(context *gin.Context) {
		upstreamResponse, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || upstreamResponse.StatusCode != http.StatusOK {
			context.Status(http.StatusServiceUnavailable)
			return
		}

		reader := upstreamResponse.Body
		contentLength := upstreamResponse.ContentLength
		contentType := upstreamResponse.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}
		context.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	engine.Run()
}

func htmlReadering() {
	engine := gin.Default()

	//engine.LoadHTMLGlob("templates/*")
	//engine.GET("/index", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "index.tmpl", gin.H{
	//		"title": "Main website",
	//	})
	//})

	//engine.LoadHTMLFiles("templates/template1.html")
	//engine.GET("/index", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "template1.html", gin.H{
	//		"title": "Main website",
	//	})
	//})

	//engine.LoadHTMLGlob("templates/**/*")
	//engine.GET("/posts/index", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
	//		"title": "Posts",
	//	})
	//})
	//engine.GET("/users/index", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
	//		"title": "Users",
	//	})
	//})

	// Using templates with same name in different directories
	engine.LoadHTMLGlob("templates/**/*")
	// must `define tamplates` in templates files
	engine.GET("/posts/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	// must `define tamplates` in templates files
	engine.GET("/users/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})

	engine.Run()
}

func customTemplateRenderer() {
	engine := gin.Default()
	//html := template.Must(template.ParseFiles("file1", "file2"))
	html := template.Must(template.ParseFiles("file1"))
	engine.SetHTMLTemplate(html)

	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "file1", nil)
	})

	engine.Run()
}

func customDelimiters() {
	engine := gin.Default()
	engine.Delims("{[{", "}]}")
	engine.LoadHTMLGlob("templates/customDelimiters.html")
	engine.GET("/delimiter", func(context *gin.Context) {
		context.HTML(http.StatusOK, "customDelimiters.html", gin.H{"title": "delimiter"})
	})
	engine.Run()
}

// custom template filters
func customTemplateFuncs() {
	engine := gin.Default()

	formatAsDate := func(t time.Time) string {
		year, month, day := t.Date()
		return fmt.Sprintf("%d%02d/%02d", year, month, day)
	}
	engine.Delims("{[{", "}]}")
	engine.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})

	engine.LoadHTMLGlob("templates/raw.tmpl")
	engine.GET("/raw", func(context *gin.Context) {
		context.HTML(http.StatusOK, "raw.tmpl", gin.H{
			"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
		})
	})
	engine.Run()
}

// context
// redirect
func redirectTest() {
	engine := gin.Default()
	engine.GET("/redirect", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
	})
	engine.Run()
}

func customMiddleware() {
	Logger := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			t := time.Now()

			// Set example variable
			c.Set("example", "12345")

			// before request
			c.Next()

			// after request
			latency := time.Since(t)
			log.Print(latency)

			// access the status we are sending
			status := c.Writer.Status()
			log.Println(status)
		}
	}

	engine := gin.New()
	engine.Use(Logger())

	engine.GET("/customMiddleware", func(context *gin.Context) {
		exmple, ok := context.MustGet("example").(string)
		log.Println(exmple, ok)
	})
	engine.Run()
}

func main() {
	//sampleDemo()
	//parameterInPath()
	//multipartOrUrlencodeForm()
	//singleFileUploadUseForm()
	//multipleFilesUseForm()
	//groupingRoutesTest()
	//middlewareTest()
	//howToWriteLogFile()
	//modelBindingTest()
	//responseContentTypeTest()
	//serveStaticFiles()
	//servingDataFromReader()
	//htmlReadering()
	//customTemplateRenderer()
	//customDelimiters()
	//customTemplateFuncs()
	//redirectTest()
	customMiddleware()
}
