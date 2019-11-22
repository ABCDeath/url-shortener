package main

import (
    "crypto/md5"
	"fmt"
    "net/http"
    "time"
	"github.com/gin-gonic/gin"
)


type UrlRequest struct {
    Url string `form:"url" json:"url" binding:"required"`
    KeepForDays int `form:"keep_for_days" json:"keep_for_days"`
}


type UrlDB struct {
    Url string
    AccessCode string
    KeepInfinitely bool
    DeleteAt time.Time
}


const HOST = "http://127.0.0.1:8000"
const URL_CODE_LEN = 6

var url_database = map[string]*UrlDB{}


func ParseMarshalError(err error) {
    return
}


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


func set_deletion_time(url_instance *UrlDB, keep_for int) bool {
    if keep_for < 0 {
        url_instance.KeepInfinitely = true
    } else if keep_for == 0 {
        url_instance.DeleteAt = time.Time{}
    } else {
        keep_hours, err := time.ParseDuration(fmt.Sprintf("%vh", keep_for * 24))
        if err != nil {
            return false
        }

        url_instance.DeleteAt = time.Now().UTC().Add(keep_hours)
    }

    return true
}


func SaveUrl(url string, keep_for int) *UrlDB {
    url_code := generate_url_code(url)
    url_database[url_code] = &UrlDB{
        Url: url,
        AccessCode: url_code,
    }

    if ok := set_deletion_time(url_database[url_code], keep_for); !ok {
        return nil
    }

    fmt.Printf("Save url <%s>: <%v> delete at %v\n", url, url_code, url_database[url_code].DeleteAt)

    return url_database[url_code]
}


func GetUrl(url_code string) (string, bool) {
    fmt.Printf("Get url: %v from %v\n", url_code, url_database)
    url_instance, ok := url_database[url_code]
    if !ok {
        return "", false
    }

    if url_instance.DeleteAt.IsZero() {
        delete(url_database, url_code)
    }

    return url_instance.Url, true
}


func UrlAddHandler(ctx *gin.Context) {
    var body UrlRequest
    if err := ctx.ShouldBindJSON(&body); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
        return
    }

    url_instance := SaveUrl(body.Url, body.KeepForDays)
    if url_instance == nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("")})
        return
    }

    response := gin.H{
        "url": fmt.Sprintf("%s/%v", HOST, url_instance.AccessCode),
        "valid_until": fmt.Sprintf("%v", url_instance.DeleteAt),
    }

    ctx.JSON(http.StatusOK, response)
}


func UrlGetHandler(ctx *gin.Context) {
    url_code := ctx.Param("url_code")[1:]

    url, ok := GetUrl(url_code)
    if !ok {
        ctx.JSON(http.StatusNotFound, gin.H{"error": url_code})
        return
    }

    ctx.Redirect(http.StatusMovedPermanently, url)
}


func main() {
    router := gin.Default()

    router.POST("/url/add", UrlAddHandler)
    router.GET("/*url_code", UrlGetHandler)

    router.Run(":8000")
}
