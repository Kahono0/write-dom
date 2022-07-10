package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"bytes"
	// "encoding/json"
)

var (
	token  = "R5952F0QQ3ErcqjuF0uDC5NaLnClvc79imG1XjV9"
	userid = "1490371"
)

type request struct {
	cookie string
}

type respStruct struct {
	Data struct {
		AssignmentsHelpneeded struct {
			Total int `json:"total"`
			PerPage int `json:"perPage"`
			Page int `json:"page"`
			Assignments []interface{} `json:"assignments"`
		} `json:"assignments_helpneeded"`
	} `json:"data"`
	Status int `json:"status"`
	Errors []interface{} `json:"errors"`
	Alerts []interface{} `json:"alerts"`
}
//func headers to map

func timeNow() string {
	t := time.Now()
	return url.QueryEscape(t.Format("2006-01-02 15:04:05"))
}

func (r *request) login() error {
	// GET https://api.writedom.com/writer/assignments/

	// Create request
	req, err := http.NewRequest("GET", "https://api.writedom.com/writer/assignments/?user_id="+userid+"&page=1&perPage=20&access_token=none&app_id=3&_token="+token+"&is_new_wd=true&local_time="+timeNow(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// Headers
	req.Header.Add("Host", "api.writedom.com")
	req.Header.Add("Origin", "https://writedom.com")
	req.Header.Add("Cookie", "production_laravel_session_api_1=eyJpdiI6ImZJZk5zT2I3RDN4bjFXZjJIbXBtRGc9PSIsInZhbHVlIjoiL0VnS1pueWhJbGs3cFBiaG1jK0gybndUTVFGWG43T3JDNkZJY2JHRlZKWCtRL2UyWlNQTlN5cjJ1QUZ5bDJHRWRzbm51ZVhwazN0OUZWSHZHZU1XZzlqenhUZjdhV0JxUGRORGVXYUNGb2RzSDJ6MmpabDNVTHo2NFMxMWVWdXYiLCJtYWMiOiI2MDRhZGMyZTViMGMzMTU3ZDIxODAwNGQwNGVkODEyODVlM2E4ZjM1MDhkNGE0YjBhZWMyNjcxZGRkNDU1M2RlIiwidGFnIjoiIn0%3D; remember_web_59ba36addc2b2f9401580f014c7f58ea4e30989d=eyJpdiI6IlFjeWJlRUtZY09jUGFxS25rSmZtVHc9PSIsInZhbHVlIjoiYXpsTmJOZVJ0S1JmM1dJb2JiN00wVUUvbnZVYUNoc3FXdFQ5SFp2TkR1aGpBTFZJQ0hMQWdOb0FSbmx0SXlqdGcyT2s1RnowRkJwMzFmdnJDRHNhWlMyYkJuN1hWMXFxS09Wb2t2KzZ1Qm9xSGJzeWVIUnF0bVpocFpISWxUbXJqTzR2Tm5FUW1PL3FvQ21UaitKQ25qYURjdlFsSVFEMDE0NXQ1WmlLR2JuRkc4OFJlRWlIajdxdTZpeVhweXdZT2lLeEVVcEtDMk1mTzZ4amRPdzQvZSsyaFRqczlrekoxVjExRG9jdVE4QnBrVlFITEpIVkQxTXpiVm5DVG5YVCIsIm1hYyI6IjE3NWIzMjQwODExYTk3MGM2MTc0M2VkZTY5OGRmMzE2ODliYjAyOWM5OWY1ZjI4NDBhYjliZjRjNzYxMDhmNjMiLCJ0YWciOiIifQ%3D%3D; __zlcprivacy=1; __zlcmid=1Aml30fc9dOPcCn; _ym_visorc=w; _ga=GA1.2.1766508788.1656869025; _gid=GA1.2.1606384748.1656869025; _ym_isad=2; _ym_d=1656869010; _ym_uid=1656869010231835314")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Safari/605.1.15")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Referer", "https://writedom.com/")

	if err := req.ParseForm(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// Fetch Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	defer resp.Body.Close()

	// Read Response Body
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	//header to map[string]interface{}
	header := resp.Header
	for k, v := range header {
		header[strings.ToLower(k)] = v
	}
	r.cookie = string(header["set-cookie"][0])
	// Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))

	return nil
}

func (r *request) getOrders() {
	req, err := http.NewRequest("GET", "https://api.writedom.com/writer/assignments/helpneeded?user_id="+userid+"&page=1&perPage=20&access_token=none&app_id=3&_token="+token+"&is_new_wd=true&local_time="+timeNow(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	req.Header.Add("Authority", "api.writedom.com")
	req.Header.Add("Method", "GET")
	// req.Header.Add("Path", "/writer/assignments/helpneeded?user_id="+userid+"&page=1&perPage=20&access_token=none&app_id=3&_token="+token+"&is_new_wd=true&local_time="+timeNow())
	req.Header.Add("Scheme", "https")
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("Origin", "https://writedom.com")
	req.Header.Add("Cookie", r.cookie)
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.9")
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Safari/605.1.15")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("origin", "https://writedom.com")
	req.Header.Add("referer", "https://writedom.com/")
	req.Header.Set("sec-fetch-dest", "\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Linux")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	if err := req.ParseForm(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer resp.Body.Close()
	//unzip response body
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Read Response Body
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	//unzip response body
	respBody, err = gzip.NewReader(bytes.NewReader(resp.Body))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	//unzip response body
	respBody, err = ioutil.ReadAll(respBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	//bytes decode
	// var data map[string]interface{}
	// err = json.Unmarshal(respBody, &data)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	return
	// }
	// fmt.Println(data)
	// Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// // fmt.Println("response Body : ", string(respBody))
	// fmt.Println(string(respBody))

	//json decode
	// var data map[string]interface{}
	// json.NewDecoder(resp.Body).Decode(&data)
	// fmt.Println(data)
	// respBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	return
	// }
	// if resp.Status != "200 OK" {
	// 	fmt.Fprintln(os.Stderr, "status code is not 200")
	// 	return
	// }
	//btes to string
	// body := string(respBody)
	// //fmt.Println(body)
	// fmt.Println(body)
// 	var data respStruct
// 	err = json.Unmarshal(respBody, &data)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		return
// 	}
// 	fmt.Println(data)
}
func main() {
	r := &request{}
	r.login()
	r.getOrders()
}


