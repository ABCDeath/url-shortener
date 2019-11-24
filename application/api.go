package main

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)


type UrlRequest struct {
    Url string `form:"url" json:"url" binding:"required"`
    KeepForDays uint `form:"keep_for_days" json:"keep_for_days"`
}


func UrlAddHandler(ctx *gin.Context) {
    var body UrlRequest
    if err := ctx.ShouldBindJSON(&body); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
        return
    }

    url_instance, err := SaveUrl(body.Url, body.KeepForDays)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
        return
    }

    response := gin.H{
        "url": fmt.Sprintf("%s/%v", HOST, url_instance.Id),
        "valid_until": fmt.Sprintf("%v", url_instance.DeleteAt),
    }

    ctx.JSON(http.StatusOK, response)
}


func UrlGetHandler(ctx *gin.Context) {
    url_code := ctx.Param("url_code")[1:]

    url, err := GetUrl(url_code)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": url_code})
        return
    }

    ctx.Redirect(http.StatusMovedPermanently, url)
}
