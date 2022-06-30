package main

import (
	"start-gin/configs"
	"start-gin/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run db
	configs.ConnectDB()

	//routes
	routes.UserRoute(router)




	// router.GET("/", func(c *gin.Context){
	// 	c.JSON(200, gin.H{
	// 		"status": 200,
	// 		"data": "hello. This is first program with GIN",
	// 	})
	// })




	router.Run("localhost:3000")
}