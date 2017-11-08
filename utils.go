package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

const instagramHost = "https://www.instagram.com/"

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

func (i *InstagramPhotos) fetchInstagramAPI(p *InstagramPage) {
	body := makeRequest(p.APIURL)
	json.Unmarshal(body, &i)
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
	splitURL := strings.Split(p.PhotoURL, "/")
	p.Body = splitHTML
	p.URLHash = splitURL[4]

	for _, t := range strings.Fields(p.Body[103]) {
		if strings.HasPrefix(t, "(@") {
			userName := strings.TrimRight(strings.TrimLeft(t, "(@"), ")")
			p.Username = userName
		}
	}

	p.APIURL = instagramHost + p.Username + "/media"

	return nil
}

func (p *InstagramPage) filterImages(instagramPhotos InstagramPhotos) {
	p.Images = p.Images[:0]

	for _, item := range instagramPhotos.Items {
		if item.Code == p.URLHash {
			if len(item.CarouselMedia) > 0 {
				for photoIndex := range item.CarouselMedia {
					p.Images = append(
						p.Images,
						item.CarouselMedia[photoIndex].Images.StandardResolution.URL,
					)
				}
			} else {
				p.Images = append(
					p.Images,
					item.Images.StandardResolution.URL,
				)
			}
		}
	}
}

func (p *InstagramPage) fetchMultiplePhotos() {
	p.BotMessage = p.BotMessage[:0]
	for i := range p.Images {
		p.BotMessage = append(
			p.BotMessage,
			linebot.NewImageMessage(p.Images[i], p.Images[i]),
		)
	}
}
