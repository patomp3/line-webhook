package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func init() {
	godotenv.Load()
}

func main() {
	var err error
	// สร้าง LINE Bot Client
	bot, err = linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal("Error initializing bot:", err)
	}

	// ใช้ Gin Framework
	r := gin.Default()

	// Health Check Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Webhook service is running!"})
	})

	// Webhook Endpoint
	r.POST("/webhook", webhookHandler)

	// Port ที่ Railway ใช้ค่าจาก Environment Variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ค่า Default
	}

	log.Println("Server started on port", port)
	r.Run(":" + port)
}

// Webhook Handler
func webhookHandler(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			// ตรวจสอบว่า event มาจาก Group หรือไม่
			if event.Source.Type == linebot.EventSourceTypeGroup {
				if message, ok := event.Message.(*linebot.TextMessage); ok {
					if message.Text == "bot" {
						replyMessage := linebot.NewTextMessage("Hello, World!")
						_, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do()
						if err != nil {
							log.Println("Error sending reply:", err)
						}
					}
				}
			} else {
				replyMessage := linebot.NewTextMessage("Hello!!!")
				_, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do()
				if err != nil {
					log.Println("Error sending reply:", err)
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
