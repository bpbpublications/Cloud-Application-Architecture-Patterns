package main
import (
 "fmt"
 "log"
 "net/http"
)
var fail = true // Simulate failure state
func main() {
 http.HandleFunc("/downstream", func(w http.ResponseWriter, r *http.Request) {
  if fail {
   http.Error(w, "Simulated downstream failure", http.StatusInternalServerError)
   return
  }
  fmt.Fprintln(w, "Success from downstream!")
 })
 fmt.Println("Downstream service running on :8081")
 log.Fatal(http.ListenAndServe(":8081", nil))
}
