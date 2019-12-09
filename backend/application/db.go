package main

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)


var client *mongo.Client


func OpenDBConnection() {
    var err error
    client, err = mongo.NewClient(options.Client().ApplyURI(DATABASE_URL))
    if err != nil {
        log.Fatalf("Error creating MongoDB client: %v\n", err)
    }

    err = client.Connect(context.TODO())
    if err != nil {
        log.Fatalf("Error connecting to MongoDB: %v\n", err)
    }

    err = client.Ping(context.TODO(), readpref.Primary())
    if err != nil {
        log.Fatalf("Error ping MongoDB: %v\n", err)
    }
}


func CloseDBConnection() {
    err := client.Disconnect(context.TODO())
    if err != nil {
        log.Fatalf("Error disconnecting from MongoDB: %v\n", err)
    }
}


func GetDBCollection(collection string) *mongo.Collection {
    return client.Database(DATABASE).Collection(collection)
}
