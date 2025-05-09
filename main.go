package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	r.GET("/test", test)

	// New Routes
	// r.POST("/sendmessage", sendMessage)
	// r.POST("/sendflexmessage", sendFlexMessage)

	// Port ที่ Railway ใช้ค่าจาก Environment Variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ค่า Default
	}

	log.Println("Server started on port", port)
	r.Run(":" + port)
}

// func sendMessage(c *gin.Context) {
// 	var req struct {
// 		UserID  string `json:"userId"`
// 		Message string `json:"message"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	msg := linebot.NewTextMessage(req.Message)
// 	_, err := bot.PushMessage(req.UserID, msg).Do()
// 	if err != nil {
// 		log.Println("Error sending message:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
// }

// func sendFlexMessage(c *gin.Context) {
// 	var req struct {
// 		UserID string `json:"userId"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	flexContainer := createFlexMessage()
// 	if flexContainer == nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create flex message"})
// 		return
// 	}

// 	msg := linebot.NewFlexMessage("สรุปยอดเงินออม", flexContainer)
// 	_, err := bot.PushMessage(req.UserID, msg).Do()
// 	if err != nil {
// 		log.Println("Error sending flex message:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "flex message sent"})
// }

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
					if message.Text == "เงินออมวันนี้" || message.Text == "เงินออม" {
						replyMessage := linebot.NewTextMessage("วันนี้ออมเงิน " + strconv.Itoa(getDay()) + " บาท")
						_, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do()
						if err != nil {
							log.Println("Error sending reply:", err)
						}
					} else if message.Text == "สรุปยอดเงินออม" || message.Text == "สรุปเงิน" {
						replyMessage := linebot.NewTextMessage("สรุปยอดเงินทั้งหมดจนถึงวันนี้คือ " + strconv.Itoa(getTotal()) + " บาท")
						_, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do()
						if err != nil {
							log.Println("Error sending reply:", err)
						}
					} else if message.Text == "สรุป" {
						flexContainer := createFlexMessage()
						replyMessage := linebot.NewFlexMessage("สรุปยอดเงินออม", flexContainer)
						_, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do()
						if err != nil {
							log.Println("Error sending Flex Message:", err)
						}
					} else if strings.HasPrefix(message.Text, "บันทึกนัดหมาย ") {
						parts := strings.SplitN(message.Text, " ", 4)
						if len(parts) == 4 {
							msg := parts[1]
							date := parts[2]
							timeStr := parts[3]

							err := saveAppointmentToMongo(event.Source.GroupID, msg, date, timeStr)
							if err != nil {
								log.Println("Error saving appointment:", err)
								reply := linebot.NewTextMessage("เกิดข้อผิดพลาดในการบันทึกนัดหมาย")
								bot.ReplyMessage(event.ReplyToken, reply).Do()
							} else {
								reply := linebot.NewTextMessage("บันทึกนัดหมายเรียบร้อยแล้ว: " + msg + " " + date + " " + timeStr)
								bot.ReplyMessage(event.ReplyToken, reply).Do()
							}
						} else {
							reply := linebot.NewTextMessage("รูปแบบคำสั่งไม่ถูกต้อง กรุณาใช้: บันทึกนัดหมาย <ข้อความ> <วันที่> <เวลา>")
							bot.ReplyMessage(event.ReplyToken, reply).Do()
						}
					}
				}
			} else {
				if message, ok := event.Message.(*linebot.TextMessage); ok {
					if message.Text == "เงินออมวันนี้" || message.Text == "aom" {
						replyMessage := linebot.NewTextMessage("วันนี้ออมเงิน " + strconv.Itoa(getDay()) + " บาท")
						_, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do()
						if err != nil {
							log.Println("Error sending reply:", err)
						}
					} else if message.Text == "สรุปยอดเงินออม" {
						replyMessage := linebot.NewTextMessage("สรุปยอดเงินทั้งหมดจนถึงวันนี้คือ " + strconv.Itoa(getTotal()) + " บาท")
						_, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do()
						if err != nil {
							log.Println("Error sending reply:", err)
						}
					} else if message.Text == "สรุป" {
						flexContainer := createFlexMessage()
						replyMessage := linebot.NewFlexMessage("สรุปยอดเงินออม", flexContainer)
						_, err := bot.ReplyMessage(event.ReplyToken, replyMessage).Do()
						if err != nil {
							log.Println("Error sending Flex Message:", err)
						}
					}
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func test(c *gin.Context) {

	saveAppointmentToMongo("11", "นัดหมาย", "9/5/2025", "13:00")
	c.JSON(http.StatusOK, gin.H{"status": getDay()})
}

func getDay() int {
	// Load the location for Thailand (Asia/Bangkok)
	// Thailand is UTC+7
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)

	// Get the current time in Thailand
	now := time.Now().In(loc)

	// Calculate the day of the year
	return now.YearDay()
}

func getTotal() int {
	// ดึงวันที่ปัจจุบัน
	//now := time.Now()

	dayNum := getDay()

	return (1 + dayNum) * dayNum / 2
}

func createFlexMessage() linebot.FlexContainer {
	flexJSON := `{
        "type": "bubble",
        "body": {
            "type": "box",
            "layout": "vertical",
            "contents": [
                {
                    "type": "text",
                    "text": "สรุปรายการประจำวัน",
                    "weight": "bold",
                    "size": "xl"
                },
				{
                    "type": "text",
                    "text": "เงินออมวันนี้: ` + strconv.Itoa(getDay()) + ` บาท",
                    "size": "md",
                    "margin": "md"
                },
                {
                    "type": "text",
                    "text": "ยอดรวม: ` + strconv.Itoa(getTotal()) + ` บาท",
                    "size": "md",
                    "margin": "md"
                }
            ]
        }
    }`

	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flexJSON))
	if err != nil {
		log.Println("Error creating Flex Message:", err)
		return nil
	}

	return flexContainer
}
