package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

func main() {
	response, err := http.Get("http://www.baidu.com")
	if err == nil {
		content, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(content))
	}

	http.Post()
}
