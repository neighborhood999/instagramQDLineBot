package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func handlerFn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World"))
}

func TestMakeRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerFn))
	defer ts.Close()

	body := makeRequest(ts.URL)
	assert.Equal(t, "Hello, World", string(body))
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
