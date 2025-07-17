package main
import (
    "context"
    "fmt"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    uri := "mongodb+srv://<username>:<password>@<cluster-url>/test?retryWrites=true&w=majority"
    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    collection := client.Database("replication_demo").Collection("logs")
    doc := map[string]interface{}{
        "timestamp": time.Now(),
        "message":   "This is a replicated log entry",
    }
    res, err := collection.InsertOne(ctx, doc)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Inserted document with ID: %v\n", res.InsertedID)
}
