package main

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)


var client *mongo.Client
var ctx, _ = context.WithTimeout(context.Background(), 10 * time.Second)


func OpenDBConnection() {
    var err error
    client, err = mongo.NewClient(options.Client().ApplyURI(DATABASE_URL))
    if err != nil {
        log.Fatal(err)
    }

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }
}


func CloseDBConnection() {
    err := client.Disconnect(ctx)
    if err != nil {
        log.Fatal(err)
    }
}


func GetDBCollection(collection string) *mongo.Collection {
    return client.Database(DATABASE).Collection(collection)
}
