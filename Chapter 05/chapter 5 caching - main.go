package main

import (
    "context"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Initialize Redis client
var rdb = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Simulated expensive function
func computeRateFromDB(zipCode string) string {
    fmt.Println("Fetching rate from expensive DB/API...")
    time.Sleep(2 * time.Second) // Simulate delay
    return "$9.99"              // Dummy rate
}

func getShippingRate(zipCode string) (string, error) {
    cacheKey := "rate:" + zipCode

    // Try to get from cache
    cached, err := rdb.Get(ctx, cacheKey).Result()
    if err == nil {
        fmt.Println("Cache hit!")
        return cached, nil
    }

    fmt.Println("Cache miss. Computing rate...")
    // Cache miss – compute it
    computedRate := computeRateFromDB(zipCode)

    // Set in Redis with 1 hour expiration
    err = rdb.Set(ctx, cacheKey, computedRate, time.Hour).Err()
    if err != nil {
        return "", err
    }

    return computedRate, nil
}

func main() {
    zip := "94107"

    rate, err := getShippingRate(zip)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Shipping rate for", zip, "is", rate)
}
