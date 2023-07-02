package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"wannabe/config"
	"wannabe/controller"
)

func main() {
	//init database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Cannot connect to database", err)
		return
	}

	//init gin
	r := gin.Default()

	var AuthC = controller.AuthController{
		Db: db,
	}

	r.GET("/hi", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"Messages": "Hi",
		})
	})
	r.POST("/v1/auth/register", AuthC.Register)
	r.POST("/v1/auth/login", AuthC.Login)

	//No route Handler
	r.NoRoute(controller.NotFoundHandler)
	r.LoadHTMLGlob("templates/*")

	r.Run(":3000")
}
