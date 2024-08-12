package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver for YugabyteDB
)

type cPool struct {
	conns   []*conn
	maxConn int
	mu      sync.Mutex
}

func NewCPool(maxConn int) (*cPool, error) {
	var mu = sync.Mutex{}
	pool := &cPool{
		mu:      mu,
		conns:   make([]*conn, 0, maxConn),
		maxConn: maxConn,
	}

	for i := 0; i < maxConn; i++ {
		pool.conns = append(pool.conns, newConn())
	}

	return pool, nil
}

func (pool *cPool) Get() (*conn, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Wait until a connection is available
	for len(pool.conns) == 0 {
		// Unlock while waiting to avoid blocking other operations
		pool.mu.Unlock()
		time.Sleep(10 * time.Millisecond) // Small sleep to avoid busy-waiting
		pool.mu.Lock()
	}

	c := pool.conns[0]
	pool.conns = pool.conns[1:]
	return c, nil
}

func (pool *cPool) Put(c *conn) {
	pool.mu.Lock()
	pool.conns = append(pool.conns, c)
	pool.mu.Unlock()
}

func benchmarkWithPool(pool *cPool, maxConnections int) time.Duration {
	startTime := time.Now()
	var wg sync.WaitGroup
	connectionFailures := 0
	successfulConnections := 0

	fmt.Printf("\nStarting benchmark with %d requests (using pooling)...\n", maxConnections)

	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()
			log.Printf("[Request %d] Attempting to acquire connection from pool...", requestID)

			conn, err := pool.Get()
			if err != nil {
				log.Printf("[Request %d] No available connections, operation skipped.", requestID)
				connectionFailures++
				return
			}
			log.Printf("[Request %d] Connection acquired. Executing query...", requestID)

			_, err = conn.db.Exec("SELECT pg_sleep(0.01);")
			if err != nil {
				if err.Error() == "EOF" {
					log.Printf("[Request %d] Error: Encountered EOF, connection lost.", requestID)
				} else {
					log.Printf("[Request %d] Error: Query execution failed: %v", requestID, err)
				}
				connectionFailures++
				pool.Put(conn)
				return
			}

			successfulConnections++
			log.Printf("[Request %d] Query executed successfully. Returning connection to pool...", requestID)
			pool.Put(conn)
		}(i + 1)
	}
	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("\nBenchmark Complete (With Pooling)\n")
	fmt.Printf("Total Requests: %d\n", maxConnections)
	fmt.Printf("Successful Connections: %d\n", successfulConnections)
	fmt.Printf("Connection Failures: %d\n", connectionFailures)
	fmt.Printf("Time Taken: %v\n", duration)
	fmt.Printf("Average Time per Request: %v\n", duration/time.Duration(maxConnections))
	return duration
}

func newConn() *conn {
	connStr := "host=localhost port=5433 user=yugabyte password=yugabyte dbname=yugabyte sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil
	}

	return &conn{db: db}
}

type conn struct {
	db *sql.DB
}

func (c *conn) Close() error {
	return c.db.Close()
}

func main() {
	maxConnections := 1000
	poolSize := 10

	fmt.Println("Initializing connection pool...")
	pool, err := NewCPool(poolSize)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	fmt.Println("Running benchmark with connection pooling...")
	benchmarkWithPool(pool, maxConnections)
}
