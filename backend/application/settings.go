package main

import (
	"fmt"
	"os"
)


var HOST = os.Getenv("SERVICE_HOST")

var DATABASE_HOST = os.Getenv("DATABASE_HOST")
var DATABASE_USER = os.Getenv("DATABASE_USER")
var DATABASE_PASSWORD = os.Getenv("DATABASE_PASSWORD")
var DATABASE = os.Getenv("DATABASE")
var DATABASE_URL = fmt.Sprintf(
	"mongodb://%s:%s@%s",
	DATABASE_USER, DATABASE_PASSWORD, DATABASE_HOST)
