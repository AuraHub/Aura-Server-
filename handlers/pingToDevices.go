package handlers

import (
	"Aura-Server/initializers"
	"context"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
)

type devicePing struct {
	DeviceId string `json:"deviceId"`
}

func Ping() {
	// Send ping to devices
	initializers.PahoConnection.Publish("ping", 0, false, "")

	// Wait 3 seconds
	threeSecondsAgo := time.Now()
	time.Sleep(3 * time.Second)

	filter := bson.D{
		{Key: "online", Value: true},
		{Key: "last_online", Value: bson.D{{Key: "$lte", Value: threeSecondsAgo}}},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "online", Value: false}}},
	}

	// Remove devices which didn't respond
	_, _ = initializers.Database.Collection("devices").UpdateMany(context.TODO(), filter, update)
	_, _ = initializers.Database.Collection("deviceTriggers").UpdateMany(context.TODO(), filter, update)
}

func ReturnedPing(c mqtt.Client, m mqtt.Message) {
	// Convert data to JSON
	var pingData devicePing

	err := json.Unmarshal(m.Payload(), &pingData)
	if err != nil {
		panic(err)
	}

	// Update status to online
	filter := bson.D{{Key: "device_id", Value: pingData.DeviceId}}
	update := bson.D{
		{
			Key:   "$set",
			Value: bson.D{{Key: "online", Value: true}, {Key: "last_online", Value: time.Now()}},
		},
	}

	_, _ = initializers.Database.Collection("devices").UpdateOne(context.TODO(), filter, update)
	_, _ = initializers.Database.Collection("deviceTriggers").UpdateOne(context.TODO(), filter, update)
}
