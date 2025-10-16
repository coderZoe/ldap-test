package main

import (
	"github.com/coderZoe/ldap-test/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/user/add", handler.AddUser)
	r.GET("/user/search", handler.SearchUser)
	r.POST("/user/modify", handler.ModifyUser)
	r.POST("/user/delete/:cn", handler.DeleteUser)
	r.POST("/enterprise/init", handler.InitEnterpriseLDAP)
	r.POST("/enterprise/clear", handler.ClearEnterpriseLDAP)
	r.GET("/enterprise/crawl", handler.CrawlEnterprise)
	r.Run(":8080")
}
