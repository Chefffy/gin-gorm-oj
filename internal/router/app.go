package router

import (
	_ "gin-gorm-oj/docs"
	"gin-gorm-oj/middlewares"
	"gin-gorm-oj/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine{
	r:=gin.Default()

	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//路由规则

	//公有方法
	//problem
	r.GET("/problem-list",service.GetProblemList)
	r.GET("/problem-detail",service.GetProblemDetail)

	//user
	r.GET("/user-detail",service.GetUserDetail)
	r.POST("/login",service.Login)
	r.POST("/send-code",service.SendCode)
	r.POST("/register",service.Register)

	r.GET("/rank-list",service.GetRankList)

	//submit
	r.GET("/submit-list",service.GetSubmitList)

	//管理员私有方法
	authAdmin := r.Group("/admin",middlewares.AuthAdminCheck())

	//problem
	authAdmin.POST("/problem-create",service.ProblemCreate)
	authAdmin.PUT("/problem-modify",service.ProblemModify)


	//categoryList
	authAdmin.GET("/category-list",service.GetCategoryList)


	//category
	authAdmin.POST("/category-create",service.CategoryCreate)
	authAdmin.PUT("/category-modify",service.CategoryModify)
	authAdmin.DELETE("/category-delete",service.CategoryDelete)

	//用户私有方法
	authUser := r.Group("/user",middlewares.AuthUserCheck())

	authUser.POST("/submit",service.Submit)

	return r
}
