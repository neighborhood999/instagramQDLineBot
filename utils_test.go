package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/stretchr/testify/assert"
)

func testCallbackHandler(w http.ResponseWriter, r *http.Request) {
	json, _ := ioutil.ReadFile("tests/media.json")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func testCallbackHandlerWithHTML(w http.ResponseWriter, r *http.Request) {
	html, _ := ioutil.ReadFile("tests/testHTML")
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func TestMakeRequest(t *testing.T) {
	expectedResponse, _ := ioutil.ReadFile("tests/media.json")
	ts := httptest.NewServer(http.HandlerFunc(testCallbackHandler))
	defer ts.Close()

	body := makeRequest(ts.URL)
	assert.Equal(t, expectedResponse, body)
}

func TestValidateURL(t *testing.T) {
	p := InstagramPage{}
	expectedHostOne := "https://www.instagram.com/"
	expectedHostTwo := "https://instagram.com"
	expectedResponseMessage := "⚠️ 請點選 Instagram 照片 [⋯] 圖示並複製網址！"
	expectedUnexpectedURLResponse := "😣 請不要輸入 Instagram 以外的網址！"

	p.validateURL(expectedHostOne)
	assert.Equal(t, expectedHostOne, p.PhotoURL)
	p.validateURL(expectedHostTwo)
	assert.Equal(t, expectedHostTwo, p.PhotoURL)
	errResponseMessage := p.validateURL("Hello LineBot")
	assert.EqualError(t, errResponseMessage, expectedResponseMessage)
	errUnexpectedURLResponse := p.validateURL("https://www.google.com.tw")
	assert.EqualError(t, errUnexpectedURLResponse, expectedUnexpectedURLResponse)
}

func TestInstagramPageContent(t *testing.T) {
	p := &InstagramPage{}
	mockLineBotTextMessage := linebot.NewTextMessage("Hello World")
	expectedValidateURLMessage := "⚠️ 請點選 Instagram 照片 [⋯] 圖示並複製網址！"

	err := p.instagramPageContent(mockLineBotTextMessage)
	assert.EqualError(t, err, expectedValidateURLMessage)

	ts := httptest.NewServer(http.HandlerFunc(testCallbackHandlerWithHTML))
	defer ts.Close()

	p.instagramPageContent(linebot.NewTextMessage(ts.URL + "/p/Ba0ExjJhvtX/"))
	assert.NotEmpty(t, p.Body)
}
