This is a simple url-shorten service written in Go, based on [gin framework](https://github.com/gin-gonic/gin) and MongoDB. To manage dependencies i use [govendor](https://github.com/kardianos/govendor).  

### API
#### `/url/add` POST
```json
{
    "url": "string",
    "keep_for_days": "int"
}
```  
##### Parameters:  
`keep_for_days` could be ommited and have to be >= 0. If `keep_for_days` == 0, the entry is deleted right after server returns it's url by id.  
##### Response:  
- `200`:
```json
{
    "url": "string",
    "valid_until": "string"
}
```  
`url` contains short url code  
`valid_until` is a datetime string `yyyy-mm-dd HH:MM:ss.f +z`  
- `400`:
```json
{
    "error": "string"
}
```  
`error` containes error description

#### `/{url_id}` GET
##### Response:
- `301` with `Location` header containing the original url  
- `404` if url_id can not be found
