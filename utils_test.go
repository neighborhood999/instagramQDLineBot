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
	json, _ := ioutil.ReadFile("media.json")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func testCallbackHandlerWithHTML(w http.ResponseWriter, r *http.Request) {
	html, _ := ioutil.ReadFile("testHTML")
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func TestMakeRequest(t *testing.T) {
	expectedResponse, _ := ioutil.ReadFile("media.json")
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

func TestFetchInstagramAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testCallbackHandler))
	defer ts.Close()

	p := &InstagramPage{PhotoURL: ts.URL}
	i := &InstagramPhotos{}
	i.fetchInstagramAPI(p)

	assert.EqualValues(t, 20, len(i.Items))
}

func TestInstagramPageContent(t *testing.T) {
	p := &InstagramPage{}
	mockLineBotTextMessage := linebot.NewTextMessage("Hello World")

	expectedValidateURLMessage := "⚠️ 請點選 Instagram 照片 [⋯] 圖示並複製網址！"
	expectedUsername := "unsplash"
	expectedURLHash := "Ba0ExjJhvtX"

	err := p.instagramPageContent(mockLineBotTextMessage)
	assert.EqualError(t, err, expectedValidateURLMessage)

	ts := httptest.NewServer(http.HandlerFunc(testCallbackHandlerWithHTML))
	defer ts.Close()

	p.instagramPageContent(linebot.NewTextMessage(ts.URL + "/p/Ba0ExjJhvtX/"))
	assert.Equal(t, expectedUsername, p.Username)
	assert.Equal(t, expectedURLHash, p.URLHash)
}

func TestFilterImages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testCallbackHandler))
	defer ts.Close()

	i := &InstagramPhotos{}
	p := &InstagramPage{
		Username: "unsplash",
		PhotoURL: ts.URL + "/p/Ba0ExjJhvtX/",
		URLHash:  "Ba0ExjJhvtX",
	}

	i.fetchInstagramAPI(p)
	p.filterImages(*i)
	assert.Equal(t, "Ba0ExjJhvtX", p.URLHash)
	assert.Equal(t, 1, len(p.Images))

	// Test multiple photos and check `p.Images` will clear
	p.PhotoURL = ts.URL + "/p/Baydi6BB99r"
	p.URLHash = "Baydi6BB99r"
	p.filterImages(*i)
	assert.Equal(t, "Baydi6BB99r", p.URLHash)
	assert.Equal(t, 10, len(p.Images))
}

func TestFetchMultiplePhotos(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testCallbackHandler))
	defer ts.Close()

	i := &InstagramPhotos{}
	p := &InstagramPage{
		Username: "unsplash",
		PhotoURL: ts.URL + "/p/Ba0ExjJhvtX/",
		URLHash:  "Ba0ExjJhvtX",
	}

	i.fetchInstagramAPI(p)
	p.filterImages(*i)
	p.fetchMultiplePhotos()
	assert.Equal(t, 1, len(p.BotMessage))

	// Test multiple photos and check `p.BotMessage` will clear
	p.PhotoURL = ts.URL + "/p/Baydi6BB99r"
	p.URLHash = "Baydi6BB99r"
	p.filterImages(*i)
	p.fetchMultiplePhotos()
	assert.Equal(t, 10, len(p.BotMessage))
}
