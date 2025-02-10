package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// ใช้ Gin Framework
	r := gin.Default()

	// Health Check Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Webhook service is running!"})
	})

	// Webhook Endpoint
	r.POST("/webhook", func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		fmt.Println("Received Webhook:", payload)
		c.JSON(http.StatusOK, gin.H{"status": "success2"})
	})

	// Port ที่ Railway ใช้ค่าจาก Environment Variable
	port := "8080"
	// if envPort := getEnv("PORT", "8080"); envPort != "" {
	// 	port = envPort
	// }

	r.Run(":" + port) // Start Server
}

// getEnv ฟังก์ชันช่วยดึงค่าตัวแปร Environment
// func getEnv(key, defaultValue string) string {
// 	if value := getenv(key); value != "" {
// 		return value
// 	}
// 	return defaultValue
// }
