package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
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

func (p *InstagramPage) validateURL(text string) error {
	var errMessage error

	url, err := url.ParseRequestURI(text)
	if err != nil {
		errMessage = errors.New("⚠️ 請點選 Instagram 照片 [⋯] 圖示並複製網址！")
		return errMessage
	}

	if url.Hostname() == "www.instagram.com" ||
		url.Hostname() == "instagram.com" ||
		url.Hostname() == "127.0.0.1" {
		p.PhotoURL = url.String()
		return nil
	}
	errMessage = errors.New("😣 請不要輸入 Instagram 以外的網址！")

	return errMessage
}

func (p *InstagramPage) instagramPageContent(text *linebot.TextMessage) error {
	if err := p.validateURL(text.Text); err != nil {
		return err
	}

	body := makeRequest(p.PhotoURL)
	splitHTML := strings.Split(string(body), "\"")
	p.Body = splitHTML

	for i := range splitHTML {
		if strings.Contains(splitHTML[i], "og:image") {
			p.PhotoURL = splitHTML[i+2]
		}
	}

	return nil
}
