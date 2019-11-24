package main

import (
    "github.com/gin-gonic/gin"
)


func main() {
    OpenDBConnection()
    defer CloseDBConnection()

    go DeleteOldUrls()

    router := gin.Default()

    router.POST("/url/add", UrlAddHandler)
    router.GET("/*url_code", UrlGetHandler)

    router.Run(":8000")
}
