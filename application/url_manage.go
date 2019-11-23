package main

import (
    "context"
    "crypto/md5"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)


type UrlDB struct {
    Id string `bson:"_id",omitempty`
    Url string
    KeepInfinitely bool `bson:"keep_infinitely"`
    DeleteAt time.Time `bson:"delete_at"`
}


const URL_CODE_LEN = 6
const URL_COLLECTION = "url"


func generate_url_code(url string, collection *mongo.Collection) string {
    url_md5 := fmt.Sprintf("%x", md5.Sum([]byte(url)))
    var url_instance UrlDB

    for i := 0; i < md5.Size * 2 - URL_CODE_LEN; i++ {
        url_id := url_md5[i:i + URL_CODE_LEN]

        err := collection.FindOne(context.TODO(), bson.D{{"_id", url_id}}).Decode(&url_instance)
        if err != nil {
            return url_id
        }

        if url_instance.Url == url {
            return url_instance.Id
        }
    }

    return ""
}


func set_deletion_time(url_instance *UrlDB, keep_for int) {
    if keep_for < 0 {
        url_instance.KeepInfinitely = true
    } else if keep_for == 0 {
        url_instance.DeleteAt = time.Time{}
    } else {
        keep_hours, err := time.ParseDuration(fmt.Sprintf("%vh", keep_for * 24))
        if err != nil {
            return
        }

        url_instance.DeleteAt = time.Now().UTC().Add(keep_hours)
    }
}


func SaveUrl(url string, keep_for int) *UrlDB {
    collection := GetDBCollection(URL_COLLECTION)

    url_code := generate_url_code(url, collection)
    url_instance := UrlDB{Id: url_code, Url: url}

    set_deletion_time(&url_instance, keep_for)

    _, err := collection.InsertOne(context.TODO(), url_instance)
    if err != nil {
        log.Fatalf("Error inserting url: %v\n", err)
    }

    return &url_instance
}


func GetUrl(url_code string) (string, bool) {
    collection := GetDBCollection(URL_COLLECTION)

    var url_instance UrlDB
    err := collection.FindOne(context.TODO(), bson.D{{"_id", url_code}}).Decode(&url_instance)
    if err != nil {
        return "", false
    }

    if url_instance.DeleteAt.IsZero() {
        _, err := collection.DeleteOne(context.TODO(), bson.D{{"_id", url_code}})
        if err != nil {
            log.Fatalf("Error deleting url: %v\n", err)
        }
    }

    return url_instance.Url, true
}
