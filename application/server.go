package main

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)


type UrlRequest struct {
    Url string `form:"url" json:"url" binding:"required"`
    KeepForDays int `form:"keep_for_days" json:"keep_for_days"`
}


const HOST = "http://127.0.0.1:8000"


func ParseMarshalError(err error) {
    // todo
    return
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
