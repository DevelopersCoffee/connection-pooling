package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver for YugabyteDB
)

func benchmarkNonPool(maxConnections int) {
	var wg sync.WaitGroup
	connectionFailures := 0
	successfulConnections := 0

	fmt.Printf("\nStarting benchmark with %d connections (no pooling)...\n", maxConnections)

	startTime := time.Now()

	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go func(connID int) {
			defer wg.Done()
			db := newConn()
			if db == nil {
				log.Printf("[Conn %d] Skipping operation due to failed connection", connID)
				connectionFailures++
				return
			}
			_, err := db.db.Exec("SELECT pg_sleep(0.01);")
			if err != nil {
				if err.Error() == "EOF" {
					log.Printf("[Conn %d] Encountered EOF, connection lost", connID)
				} else {
					log.Printf("[Conn %d] Query execution failed: %v", connID, err)
				}
				connectionFailures++
				db.Close()
				return
			}
			successfulConnections++
			log.Printf("[Conn %d] Query executed successfully.", connID)
			db.Close()
		}(i + 1)
	}
	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("\nBenchmark Complete (No Pooling)\n")
	fmt.Printf("Total Requests: %d\n", maxConnections)
	fmt.Printf("Successful Connections: %d\n", successfulConnections)
	fmt.Printf("Connection Failures: %d\n", connectionFailures)
	fmt.Printf("Time Taken: %v\n", duration)
	fmt.Printf("Average Time per Request: %v\n", duration/time.Duration(maxConnections))
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
	fmt.Println("Running benchmark without connection pooling...")
	benchmarkNonPool(maxConnections)
}
