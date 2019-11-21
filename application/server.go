package main

import (
	"net/http"
	"fmt"
    "strconv"
	"github.com/gin-gonic/gin"
)


type UrlRequest struct {
    Url string `form:"url" json:"url" binding:"required"`
    AliveDays int `form:"alive_days" json:"alive_days"`
}


type UrlDB struct {
    Url string
    AliveDays int
    // todo: replace with timestamp or something
    DeleteAt int
}


const HOST = "http://127.0.0.1:8000"

var url_id_count = uint64(1)
var url_database = map[uint64]UrlDB{}


func ParseMarshalError(err error) {
    return
}


func SaveUrl(url string, alive int) uint64 {
    url_id := url_id_count
    url_database[url_id] = UrlDB{Url: url, AliveDays: alive}
    url_id_count++

    fmt.Printf("Save url <%s> for %v days: <%v>\n", url, alive, url_id)

    return url_id
}


func GetUrl(url_id uint64) (string, bool) {
    fmt.Printf("Get url: %v from %v\n", url_id, url_database)
    url_instance, ok := url_database[url_id]
    if !ok {
        return "", false
    }

    if url_instance.AliveDays == 0 {
        delete(url_database, url_id)
    }

    return url_instance.Url, true
}


func UrlAddHandler(ctx *gin.Context) {
    var body UrlRequest
    if err := ctx.ShouldBindJSON(&body); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
        return
    }

    url_id := SaveUrl(body.Url, body.AliveDays)

    response := gin.H{
        "url": fmt.Sprintf("%s/%v", HOST, url_id),
        "alive": fmt.Sprintf("%v days", body.AliveDays),
    }

    ctx.JSON(http.StatusOK, response)
}


func UrlGetHandler(ctx *gin.Context) {
    url_id, err := strconv.ParseUint(ctx.Param("url_id")[1:], 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": url_id})
        return
    }

    url, ok := GetUrl(url_id)
    if !ok {
        ctx.JSON(http.StatusNotFound, gin.H{"error": url_id})
        return
    }

    ctx.Redirect(http.StatusMovedPermanently, url)
}


func main() {
    fmt.Println("running...")

    router := gin.Default()

    router.POST("/url/add", UrlAddHandler)
    router.GET("/*url_id", UrlGetHandler)

    router.Run(":8000")
}
