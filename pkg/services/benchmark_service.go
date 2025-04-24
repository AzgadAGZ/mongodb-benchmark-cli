package services

import (
	"context"
	"fmt"
	"mongodb-benchmark/pkg/config"
	"mongodb-benchmark/pkg/models"
	"mongodb-benchmark/pkg/repository"
	"time"
)

type BenchmarkService struct {
	repo *repository.MongoRepository
	cfg  *config.BenchmarkConfig
}

func NewBenchmarkService(cfg *config.BenchmarkConfig) (*BenchmarkService, error) {
	// Extract database name from MongoDB URI if present
	repo, err := repository.NewMongoRepository(cfg.MongoURI, cfg.DatabaseName, cfg.Collection)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %v", err)
	}

	return &BenchmarkService{
		repo: repo,
		cfg:  cfg,
	}, nil
}

func (s *BenchmarkService) RunBenchmark() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if s.cfg.DropDb {
		fmt.Println("Attempting to drop database...")
		if err := s.repo.DropDatabase(ctx); err != nil {
			fmt.Printf("Warning: Failed to drop database: %v\nContinuing with benchmark...\n", err)
			// Don't return error, continue with benchmark
		}
	}

	if s.cfg.RunAll {
		return s.runAllTests(ctx)
	}

	return s.runSingleTest(ctx)
}

func (s *BenchmarkService) runAllTests(ctx context.Context) error {
	tests := []string{"insert", "update", "delete"}
	for _, test := range tests {
		fmt.Printf("\nRunning %s benchmark...\n", test)
		result := s.executeTest(ctx, test)
		s.printResult(result)
	}
	return nil
}

func (s *BenchmarkService) runSingleTest(ctx context.Context) error {
	fmt.Printf("\nRunning %s benchmark...\n", s.cfg.TestType)
	result := s.executeTest(ctx, s.cfg.TestType)
	s.printResult(result)
	return nil
}

func (s *BenchmarkService) executeTest(ctx context.Context, testType string) *models.BenchmarkResult {
	var result *models.BenchmarkResult

	switch testType {
	case "insert":
		result = s.repo.InsertBenchmark(ctx, s.cfg.Threads, s.cfg.Operations, s.cfg.LargeDocs)
	case "update":
		result = s.repo.UpdateBenchmark(ctx, s.cfg.Threads, s.cfg.Operations)
	case "delete":
		result = s.repo.DeleteBenchmark(ctx, s.cfg.Threads, s.cfg.Operations)
	default:
		return &models.BenchmarkResult{
			OperationType: "unknown",
			ElapsedTime:   0,
		}
	}

	return result
}

func (s *BenchmarkService) printResult(result *models.BenchmarkResult) {
	fmt.Printf("\nBenchmark Results for %s:\n", result.OperationType)
	fmt.Printf("Total Operations: %d\n", result.TotalOperations)
	fmt.Printf("Total Time: %v\n", result.ElapsedTime.Round(time.Millisecond))
	fmt.Printf("Operations/sec: %.2f\n", result.OpsPerSecond)
	fmt.Printf("Average Latency: %v\n", result.AverageLatency.Round(time.Microsecond))
}

func (s *BenchmarkService) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}
