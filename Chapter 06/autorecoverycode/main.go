package main
import (
    "fmt"
    "log"
    "net/http"
    "sync/atomic"
    "time"
)
var isHealthy int32 = 1
func main() {
    // Simulated app handler
    http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
        if atomic.LoadInt32(&isHealthy) == 0 {
            http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
            return
        }
        fmt.Fprintln(w, "Payment processed successfully!")
    })
    // Liveness probe
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if atomic.LoadInt32(&isHealthy) == 0 {
            http.Error(w, "Not healthy", http.StatusInternalServerError)
        } else {
            fmt.Fprintln(w, "Healthy")
        }
    })
    // Simulate failure after 20 seconds
    go func() {
        time.Sleep(20 * time.Second)
        log.Println("Simulating service failure...")
        atomic.StoreInt32(&isHealthy, 0)
    }()
    log.Println("Service starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
