package main

import (
    "crypto/md5"
    "fmt"
    "log"
    "time"
)


type UrlDB struct {
    Url string
    AccessCode string
    KeepInfinitely bool
    DeleteAt time.Time
}


const URL_CODE_LEN = 6

var url_database = map[string]*UrlDB{}


func generate_url_code(url string) string {
    for {
        url_md5 := fmt.Sprintf("%x", md5.Sum([]byte(url)))
        for i := 0; i < md5.Size * 2 - URL_CODE_LEN; i++ {
            url_id := url_md5[i:i + URL_CODE_LEN]
            if _, ok := url_database[url_id]; !ok {
                return url_id
            }
        }
    }
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
    url_code := generate_url_code(url)
    url_database[url_code] = &UrlDB{
        Url: url,
        AccessCode: url_code,
    }

    set_deletion_time(url_database[url_code], keep_for)

    log.Printf("Save url <%s>: <%v> delete at %v\n", url, url_code, url_database[url_code].DeleteAt)

    return url_database[url_code]
}


func GetUrl(url_code string) (string, bool) {
    log.Printf("Get url: %v from %v\n", url_code, url_database)
    url_instance, ok := url_database[url_code]
    if !ok {
        return "", false
    }

    if url_instance.DeleteAt.IsZero() {
        delete(url_database, url_code)
    }

    return url_instance.Url, true
}
