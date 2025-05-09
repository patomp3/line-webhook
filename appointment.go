package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Appointment struct {
	GroupID    string `bson:"groupId"`
	CreateDate string `bson:"createDate"`
	Message    string `bson:"message"`
	ApDate     string `bson:"apDate"`
	ApTime     string `bson:"apTime"`
}

func saveAppointmentToMongo(groupId, msg, dateStr, timeStr string) error {
	// สร้างข้อมูล
	now := time.Now().Format("02/01/2006 15:04:05")
	appointment := Appointment{
		GroupID:    groupId,
		CreateDate: now,
		Message:    msg,
		ApDate:     dateStr,
		ApTime:     timeStr,
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://mongo:QhHhaLDEsDhFbMPvBAFgwghrowhKoGgN@shortline.proxy.rlwy.net:44021")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("test").Collection("my-note")
	_, err = collection.InsertOne(context.Background(), appointment)
	return err
}
