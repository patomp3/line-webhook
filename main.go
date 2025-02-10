package main

import (
	"io"
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
	//log.Print(os.Getenv("LINE_CHANNEL_SECRET"))

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
	signature := c.GetHeader("X-Line-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing signature"})
		return
	}

	_, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
		return
	}

	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			var replyMessage string

			switch event.Source.Type {
			case linebot.EventSourceTypeGroup: // ถ้ามาจากกลุ่ม
				replyMessage = "Hello, World!"
			case linebot.EventSourceTypeUser: // ถ้ามาจาก Friend (User)
				replyMessage = "bot"
			default:
				replyMessage = "Unknown source"
			}

			_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
			if err != nil {
				log.Println("Error sending reply:", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
