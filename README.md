
# Connection Pooling vs Non-Pooling Benchmark

This repository contains two benchmarks to demonstrate the difference in performance between using a connection pool and not using one when connecting to a YugabyteDB instance.

## Repository Structure

```bash
connection-pooling/
├── connection-non-pool/
│   ├── benchmark_non_pool.go        # Benchmark code without connection pooling
│   ├── Docker-compose.yaml          # Docker Compose file to spin up YugabyteDB
│   ├── go.mod                       # Go module file
│   ├── go.sum                       # Go dependencies file
│   └── vendor/                      # Vendor directory for dependencies
├── connection-with-pool/
│   ├── benchmark_with_pool.go       # Benchmark code with connection pooling
│   ├── Docker-compose.yaml          # Docker Compose file to spin up YugabyteDB
│   ├── go.mod                       # Go module file
│   ├── go.sum                       # Go dependencies file
│   └── vendor/                      # Vendor directory for dependencies
└── README.md                        # This README file
```

## Prerequisites

- Go 1.16 or later installed
- Docker installed and running
- pgcli installed (optional, for connecting to YugabyteDB via CLI)

## Getting Started

### Step 1: Start YugabyteDB with Docker Compose

```bash
cd connection-non-pool # or connection-with-pool
docker-compose up -d
```

This command will start the YugabyteDB instance in the background.

### Step 2: Run the Benchmarks

#### Without Connection Pooling

Navigate to the `connection-non-pool` directory and run the benchmark:

```bash
cd connection-non-pool
go mod tidy
go run benchmark_non_pool.go
```

#### With Connection Pooling

Navigate to the `connection-with-pool` directory and run the benchmark:

```bash
cd connection-with-pool
go mod tidy
go run benchmark_with_pool.go
```

### Step 3: Review the Results

Compare the output of both benchmarks. You should see a significant difference in the number of successful connections and the time taken to complete all requests.

#### Example Output

##### Without Connection Pooling:

```plaintext
Benchmark Complete (No Pooling)
Total Requests: 1000
Successful Connections: 170
Connection Failures: 823
Time Taken: 1.0627465s
Average Time per Request: 1.0627465ms
```

##### With Connection Pooling:

```plaintext
Benchmark Complete (With Pooling)
Total Requests: 1000
Successful Connections: 999
Connection Failures: 0
Time Taken: 2.142920917s
Average Time per Request: 2.14292ms
```

### Step 4: Tear Down the Environment

To stop and remove the YugabyteDB containers, run:

```bash
docker-compose down
```

## Conclusion

This repository demonstrates the performance benefits of using connection pooling when working with a database like YugabyteDB. Without pooling, the application struggles to maintain a high number of concurrent connections, resulting in failures and longer execution times. With pooling, the application can efficiently reuse connections, leading to higher success rates and faster execution.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
