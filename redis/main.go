package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"log"
	"time"
)

func main() {
	var db *redis.Client
	var err error
	ctx := context.Background()
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "alpine", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})

		return db.Ping(ctx).Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	fmt.Println("connected!")
	time.Sleep(1 * time.Second)

	writer := simpleWriter{}

	opts := dockertest.ExecOptions{
		//StdIn:  writer,
		StdOut: writer,
		StdErr: writer,
	}
	exitCode, err := resource.Exec([]string{
		"redis-cli GET key1",
	}, opts)
	if err != nil {
		fmt.Println("err exec", err)
	}
	fmt.Println("exec code", exitCode)

	time.Sleep(5 * time.Second)

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
		return
	}
	fmt.Println("purged")
}

type simpleWriter struct {
}

func (s simpleWriter) Read(p []byte) (n int, err error) {
	fmt.Println("read:", string(p))
	return len(p), nil
}

func (s simpleWriter) Write(p []byte) (n int, err error) {
	fmt.Println("write:", string(p))
	return len(p), nil
}
