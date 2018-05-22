package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"mime/multipart"
	"bytes"
)

func main() {
	fileUpload()
}

func fileUpload() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("sandbox_type", "win7_sp1_enx86_office2013_sp1")
	writer.WriteField("source", "api.threatbook")
	writer.WriteField("apikey", "03c770882e87585fea0272a8e6a7b7e37085e193475884b1316e14fb193e992d")
	writer.WriteField("run_time", "60")
	writer.CreateFormFile("file", "/path/to/file")
	writer.Close()

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://s.threatbook.cn/api/v2/file/upload", body)

	// Headers
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}
