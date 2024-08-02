package routers

import (
	"userApiTest/controllers"

	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {

	r.GET("/user", controllers.UserShow)
	r.POST("/user", controllers.UserLogin)
	r.POST("/user/signup", controllers.UserSignup)
	r.PATCH("/user", controllers.UserEdit)

}
