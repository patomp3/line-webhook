package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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

func getUpcomingAppointments(groupID string) ([]Appointment, error) {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:QhHhaLDEsDhFbMPvBAFgwghrowhKoGgN@shortline.proxy.rlwy.net:44021")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("test").Collection("my-note")

	// แปลงวันที่ปัจจุบันเป็น time.Time
	today := time.Now()
	layout := "02/01/2006"

	cursor, err := collection.Find(context.Background(), bson.M{
		"groupId": groupID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var allAppointments []Appointment
	for cursor.Next(context.Background()) {
		var ap Appointment
		if err := cursor.Decode(&ap); err != nil {
			continue
		}
		// ตรวจสอบว่า apDate >= วันนี้
		apDateParsed, err := time.Parse(layout, ap.ApDate)
		if err == nil && !apDateParsed.Before(today) {
			allAppointments = append(allAppointments, ap)
		}
	}

	return allAppointments, nil
}
