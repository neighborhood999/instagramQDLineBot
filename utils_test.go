package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	json, _ := ioutil.ReadFile("media.json")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func TestMakeRequest(t *testing.T) {
	expectedResponse, _ := ioutil.ReadFile("media.json")
	ts := httptest.NewServer(http.HandlerFunc(callbackHandler))
	defer ts.Close()

	body := makeRequest(ts.URL)
	assert.Equal(t, expectedResponse, body)
}

func TestValidateURL(t *testing.T) {
	expectedHostOne := "https://www.instagram.com/"
	expectedHostTwo := "https://instagram.com"
	expectedResponseMessage := "⚠️ 請點選 Instagram 照片 [⋯] 圖示並複製網址！"
	expectedUnexpectedURLResponse := "😣 請不要輸入 Instagram 以外的網址！"

	urlOne, _ := validateURL(expectedHostOne)
	assert.Equal(t, expectedHostOne, urlOne)
	urlTwo, _ := validateURL(expectedHostTwo)
	assert.Equal(t, expectedHostTwo, urlTwo)
	_, errResponseMessage := validateURL("Hello LineBot")
	assert.EqualError(t, errResponseMessage, expectedResponseMessage)
	_, errUnexpectedURLResponse := validateURL("https://www.google.com.tw")
	assert.EqualError(t, errUnexpectedURLResponse, expectedUnexpectedURLResponse)
}

func TestFetchInstagramAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(callbackHandler))
	defer ts.Close()

	p := &ParsePage{PhotoURL: ts.URL}
	i := &InstagramPhotos{}
	i.fetchInstagramAPI(p)

	assert.EqualValues(t, 20, len(i.Items))
}
