package data

import "go.mongodb.org/mongo-driver/v2/mongo"

var client *mongo.Client


func New(mongoClient *mongo.Client) Models {
	client = mongoClient

	return Models{
		LogEntry: LogEntry{},
	}
}
type Models struct {
	LogEntry LogEntry
}
