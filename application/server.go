package main

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)


func get_cors_config_middleware() cors.Config {
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{HOST}

    return config
}


func main() {
    OpenDBConnection()
    defer CloseDBConnection()

    go DeleteOldUrls()

    router := gin.Default()

    router.Use(cors.New(get_cors_config_middleware()))

    router.POST("/url/add", UrlAddHandler)
    router.GET("/*url_code", UrlGetHandler)

    router.Run(":8000")
}
