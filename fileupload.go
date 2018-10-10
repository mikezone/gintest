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
	writer.WriteField("a", "...")
	writer.WriteField("b", "...")
	writer.WriteField("c", "...")
	writer.WriteField("d", "..")
	writer.CreateFormFile("e", "...")
	writer.Close()

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "...", body)

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
