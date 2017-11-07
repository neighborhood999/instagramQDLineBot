package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	bot             *linebot.Client
	instagramPage   InstagramPage
	instagramPhotos InstagramPhotos
	imageMessages   []linebot.Message
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

				instagramPhotos.fetchInstagramAPI(&instagramPage)
				instagramPage.filterImages(instagramPhotos)
				instagramPage.fetchMultiplePhotos()

				botMessageLength := len(instagramPage.BotMessage)
				log.Println("botMessageLength:", botMessageLength)

				if botMessageLength < 5 {
					if _, err := bot.ReplyMessage(event.ReplyToken, instagramPage.BotMessage[:botMessageLength]...).Do(); err != nil {
						log.Print(err)
					}
				} else {
					splitBotMessageMaxSend := instagramPage.BotMessage[:5]
					lastBotMessage := instagramPage.BotMessage[5:]

					if _, err := bot.ReplyMessage(event.ReplyToken, splitBotMessageMaxSend...).Do(); err != nil {
						log.Println(err)
					}
					if _, err := bot.PushMessage(event.Source.UserID, lastBotMessage...).Do(); err != nil {
						log.Println(err)
					}
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
