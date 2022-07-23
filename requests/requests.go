package requests

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

//
// DoRequest
//  @Description: 生成请求
//  @param url : 请求url
//  @param method : 请求方式 GET\POST
//  @param data : form body
//  @return response
//
func DoRequest(url string, method string, data map[string]string) (response string, err error) {

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	for key, value := range data {
		_ = writer.WriteField(key, value)
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
