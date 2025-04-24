package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Value     int                `bson:"value"`
	Data      string             `bson:"data"`
	Timestamp time.Time          `bson:"timestamp"`
}

type BenchmarkResult struct {
	OperationType   string
	TotalOperations int64
	ElapsedTime     time.Duration
	OpsPerSecond    float64
	AverageLatency  time.Duration
}

func NewDocument(large bool) Document {
	doc := Document{
		Name:      "test-document",
		Value:     1000,
		Timestamp: time.Now(),
	}

	if large {
		// Create a 1KB string for large documents
		data := make([]byte, 1024)
		for i := range data {
			data[i] = 'x'
		}
		doc.Data = string(data)
	}

	return doc
}
