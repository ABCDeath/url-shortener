package main

import (
    "crypto/md5"
	"fmt"
    "net/http"
	"github.com/gin-gonic/gin"
)


type UrlRequest struct {
    Url string `form:"url" json:"url" binding:"required"`
    AliveDays int `form:"alive_days" json:"alive_days"`
}


type UrlDB struct {
    Url string
    AccessCode string
    AliveDays int
    // todo: replace with timestamp or something
    DeleteAt int
}


const HOST = "http://127.0.0.1:8000"
const URL_CODE_LEN = 6

var url_database = map[string]UrlDB{}


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


func SaveUrl(url string, alive int) string {
    url_code := generate_url_code(url)
    url_database[url_code] = UrlDB{
        Url: url,
        AccessCode: url_code,
        AliveDays: alive,
    }

    fmt.Printf("Save url <%s> for %v days: <%v>\n", url, alive, url_code)

    return url_code
}


func GetUrl(url_code string) (string, bool) {
    fmt.Printf("Get url: %v from %v\n", url_code, url_database)
    url_instance, ok := url_database[url_code]
    if !ok {
        return "", false
    }

    if url_instance.AliveDays == 0 {
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

    url_code := SaveUrl(body.Url, body.AliveDays)

    response := gin.H{
        "url": fmt.Sprintf("%s/%v", HOST, url_code),
        "alive": fmt.Sprintf("%v days", body.AliveDays),
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
    fmt.Println("running...")

    router := gin.Default()

    router.POST("/url/add", UrlAddHandler)
    router.GET("/*url_code", UrlGetHandler)

    router.Run(":8000")
}
