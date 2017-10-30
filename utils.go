package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func makeRequest(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func validateURL(text string) (string, error) {
	url, err := url.ParseRequestURI(text)
	if err != nil {
		errMessage := errors.New("⚠️ 請點選 Instagram 照片 [⋯] 圖示並複製網址！")
		return "", errMessage
	}

	if url.Hostname() == "www.instagram.com" || url.Hostname() == "instagram.com" {
		return url.String(), nil
	}

	errMessage := errors.New("😣 請不要輸入 Instagram 以外的網址！")
	return "", errMessage
}
