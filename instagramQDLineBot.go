package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	bot           *linebot.Client
	instagramPage InstagramPage
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Println("\nUser:", event.Source.UserID, "\nIncoming Message:", message.Text)

				if validateURLError := instagramPage.instagramPageContent(message); validateURLError != nil {
					if _, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(validateURLError.Error()),
					).Do(); err != nil {
						log.Println(err)
					}
					return
				}

				if _, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewImageMessage(instagramPage.PhotoURL, instagramPage.PhotoURL),
				).Do(); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func main() {
	var err error

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", callbackHandler)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
