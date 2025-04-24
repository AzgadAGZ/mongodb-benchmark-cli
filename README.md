# MongoDB Benchmark Cli

A high-performance MongoDB benchmarking tool written in Go that tests various operations (insert, update, delete) using concurrent goroutines.

## Features

- Concurrent operations using configurable number of threads
- Support for multiple operation types (insert, update, delete)
- Option to run all tests sequentially
- Support for large document testing
- Detailed performance metrics (operations/second, average latency)
- Graceful shutdown handling

## Prerequisites

- Go 1.16 or later
- MongoDB instance (local or remote)

## Installation

```bash
git clone <repository-url>
cd mongodb-benchmark
go mod tidy
```

## Usage

Build and run the benchmark tool:

```bash
go build -o benchmark cmd/benchmark/main.go
./benchmark [flags]
```

### Available Flags

- `-uri string`: MongoDB URI (default "mongodb://localhost:27017")
- `-threads int`: Number of concurrent threads (default 10)
- `-ops int`: Number of operations to perform (default 1000)
- `-type string`: Test type (insert, update, delete) (default "insert")
- `-all`: Run all test types
- `-large`: Use large documents
- `-dropDb`: Drop database before running tests (default true)
- `-db string`: Database name (default "benchmark")
- `-coll string`: Collection name (default "testdata")

### Examples

1. Run insert benchmark with 20 threads and 10000 operations:
```bash
./benchmark -threads 20 -ops 10000 -type insert
```

2. Run all tests with large documents:
```bash
./benchmark -all -large
```

3. Run update benchmark with custom database and collection:
```bash
./benchmark -type update -db mydb -coll mycoll
```

## Output

The tool provides detailed output including:
- Configuration settings
- Progress updates
- Final results with:
  - Total operations completed
  - Total time taken
  - Operations per second
  - Average latency

## Architecture

The application follows clean architecture principles with the following structure:

```
.
├── cmd/
│   └── benchmark/
│       └── main.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── document.go
│   ├── repository/
│   │   └── mongodb.go
│   └── services/
│       └── benchmark_service.go
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 