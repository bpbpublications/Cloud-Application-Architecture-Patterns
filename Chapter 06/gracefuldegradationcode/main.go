package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type Product struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price string `json:"price"`
}

// Static fallback list
var fallbackProducts = []Product{
    {ID: 1, Name: "Wireless Mouse", Price: "$25"},
    {ID: 2, Name: "Mechanical Keyboard", Price: "$70"},
    {ID: 3, Name: "USB-C Hub", Price: "$40"},
}

// Calls the recommendation service with timeout
func fetchRecommendations() ([]Product, error) {
    client := http.Client{Timeout: 2 * time.Second}
    resp, err := client.Get("http://localhost:9090/recommend")
    if err != nil {
        return fallbackProducts, err // graceful fallback
    }
    defer resp.Body.Close()

    var products []Product
    if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
        return fallbackProducts, err
    }
    return products, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    products, err := fetchRecommendations()
    if err != nil {
        log.Println("Using fallback due to error:", err)
    }
    json.NewEncoder(w).Encode(products)
}

func main() {
    http.HandleFunc("/products", handler)
    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
