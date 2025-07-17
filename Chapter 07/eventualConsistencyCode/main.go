package main
import (
 "context"
 "fmt"
 "time"
 "go.mongodb.org/mongo-driver/bson"
 "go.mongodb.org/mongo-driver/mongo"
 "go.mongodb.org/mongo-driver/mongo/options"
)
type Shipment struct {
 ID        string    `bson:"_id"`
 Status    string    `bson:"status"`
 UpdatedAt time.Time `bson:"updated_at"`
}
func main() {
 ctx := context.TODO()
 client, err := mongo.Connect(ctx, options.Client().ApplyURI("<Your MongoDB Atlas URI>"))
 if err != nil {
  panic(err)
 }
 defer client.Disconnect(ctx)
 collection := client.Database("shipping").Collection("shipments")
 // Simulated concurrent update
 updatedStatus := Shipment{
  ID:        "SHIP123",
  Status:    "Delivered",
  UpdatedAt: time.Now(), // each service sets its own timestamp
 }
 // Use Update with $max to implement Last-Write-Wins (based on timestamp)
 filter := bson.M{"_id": updatedStatus.ID}
 update := bson.M{
  "$max": bson.M{
   "updated_at": updatedStatus.UpdatedAt,
   "status":     updatedStatus.Status,
  },
 }
 result, err := collection.UpdateOne(ctx, filter, update)
 if err != nil {
  panic(err)
 }
 fmt.Printf("Matched %v, Modified %v\n", result.MatchedCount, result.ModifiedCount)
}
