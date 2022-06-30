package routes

import (
	"github.com/gin-gonic/gin"
	"start-gin/controllers"
)


func UserRoute (router *gin.Engine){
	router.POST("/user", controllers.CreateUser())
	router.GET("/user/:userId", controllers.GetAnUser())
	router.GET("/users", controllers.GetAllUser())
	router.PATCH("/user/:userId", controllers.UpdateAnUser())
	router.DELETE("/user/:userId", controllers.DeleteAnUser())

}