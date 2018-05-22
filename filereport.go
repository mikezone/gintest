package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fileReport()
}

func fileReport() {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "https://s.threatbook.cn/api/v2/file/report?apikey=xxxx&sandbox_type=win7_sp1_enx86_office2013&sha256=0547a5e78918dd756b850798557b5cdf7ff57615cc7e6707ce3b82cd2a38bc69&source=api.threatbook&query_fileds=ioc&query_fileds=system", nil)

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
	  fmt.Println(parseFormErr)
	}

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



