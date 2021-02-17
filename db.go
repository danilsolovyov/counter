package main

import (
    "context"
    "log"
    "time"

    //"fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Counter - структура счетчика
type Counter struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    Count int             `bson:"count"`
}

// Connection - Текущее подключение к MongoDb
type Connection struct {
    *mongo.Client
}

// ConnectToDb - Подключаемся к MongoDB
func ConnectToDb() (*Connection, error) {
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    return &Connection{client}, err
}

// CreateScore - создает counter.count если он не существует
func (c *Connection) CreateScore() {
    collection := c.Database("homework").Collection("counter")
    // Проверяем существование counter.count
    x, _ := collection.CountDocuments(context.TODO(), bson.D{})
    if x == 0 {
        collection.InsertOne(context.TODO(), Counter{Count: 0})
    }
}

// GetScore - возвращает текущее значение counter.count
func (c *Connection) GetScore() Counter {
    c.CreateScore()
    collection := c.Database("homework").Collection("counter")
    var result Counter
    err := collection.FindOne(context.TODO(), bson.M{}).Decode(&result)

    if err != nil {
        log.Fatal(err)
    }
    return result
}

// UpdateScore - устанавливает новое значение counter.count
func (c *Connection) UpdateScore() {
    collection := c.Database("homework").Collection("counter")
    _, err := collection.UpdateOne(context.TODO(),
        bson.M{},
        bson.D{
            {"$inc", bson.D{{"count", 1}}},
        },
        options.Update().SetUpsert(true))

    if err != nil {
        log.Fatal(err)
    }
}
