package repository

import (
	"context"
	"crypto/tls"
	"mongodb-benchmark/pkg/models"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRepository(uri, dbName, collName string) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Configure TLS
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}

	// Configure client options for Firestore with MongoDB compatibility
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetTLSConfig(tlsConfig).
		SetTimeout(30 * time.Second).
		SetConnectTimeout(30 * time.Second).
		SetServerSelectionTimeout(30 * time.Second).
		SetDirect(false).                // Use direct connection
		SetRetryWrites(false).           // Disable retry writes for Firestore compatibility
		SetRetryReads(false).            // Disable retry reads for Firestore compatibility
		SetCompressors([]string{"none"}) // Disable compression for compatibility

	// Create client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()

	err = client.Ping(pingCtx, nil)
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	collection := client.Database(dbName).Collection(collName)
	return &MongoRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *MongoRepository) InsertBenchmark(ctx context.Context, threads, operations int, largeDocs bool) *models.BenchmarkResult {
	var wg sync.WaitGroup
	opsPerThread := operations / threads
	startTime := time.Now()

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerThread; j++ {
				doc := models.NewDocument(largeDocs)
				_, _ = r.collection.InsertOne(ctx, doc)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(startTime)
	totalOps := int64(operations)

	return &models.BenchmarkResult{
		OperationType:   "insert",
		TotalOperations: totalOps,
		ElapsedTime:     elapsed,
		OpsPerSecond:    float64(totalOps) / elapsed.Seconds(),
		AverageLatency:  elapsed / time.Duration(totalOps),
	}
}

func (r *MongoRepository) UpdateBenchmark(ctx context.Context, threads, operations int) *models.BenchmarkResult {
	var wg sync.WaitGroup
	opsPerThread := operations / threads
	startTime := time.Now()

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerThread; j++ {
				filter := bson.M{"value": 1000}
				update := bson.M{"$set": bson.M{"value": 2000}}
				_, _ = r.collection.UpdateOne(ctx, filter, update)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(startTime)
	totalOps := int64(operations)

	return &models.BenchmarkResult{
		OperationType:   "update",
		TotalOperations: totalOps,
		ElapsedTime:     elapsed,
		OpsPerSecond:    float64(totalOps) / elapsed.Seconds(),
		AverageLatency:  elapsed / time.Duration(totalOps),
	}
}

func (r *MongoRepository) DeleteBenchmark(ctx context.Context, threads, operations int) *models.BenchmarkResult {
	var wg sync.WaitGroup
	opsPerThread := operations / threads
	startTime := time.Now()

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerThread; j++ {
				filter := bson.M{"value": 1000}
				_, _ = r.collection.DeleteOne(ctx, filter)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(startTime)
	totalOps := int64(operations)

	return &models.BenchmarkResult{
		OperationType:   "delete",
		TotalOperations: totalOps,
		ElapsedTime:     elapsed,
		OpsPerSecond:    float64(totalOps) / elapsed.Seconds(),
		AverageLatency:  elapsed / time.Duration(totalOps),
	}
}

func (r *MongoRepository) DropDatabase(ctx context.Context) error {
	return r.collection.Database().Drop(ctx)
}

func (r *MongoRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
