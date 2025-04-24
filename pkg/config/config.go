package config

type BenchmarkConfig struct {
	MongoURI     string
	Threads      int
	Operations   int
	TestType     string
	Duration     int
	RunAll       bool
	LargeDocs    bool
	DropDb       bool
	DatabaseName string
	Collection   string
}

func NewDefaultConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		MongoURI:     "mongodb://localhost:27017",
		Threads:      10,
		Operations:   1000,
		TestType:     "insert",
		Duration:     0,
		RunAll:       false,
		LargeDocs:    false,
		DropDb:       true,
		DatabaseName: "benchmark",
		Collection:   "testdata",
	}
}
