package main

import (
	"userApiTest/database"
	"userApiTest/routers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	database.InitDb()
	routers.Routers(router)
	router.Run(":8080")

}
