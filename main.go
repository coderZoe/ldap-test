package main

import (
	"github.com/coderZoe/ldap-test/handler/ldap_server"
	"github.com/coderZoe/ldap-test/handler/win_ad"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// LDAP Server 路由
	r.POST("/ldapserver/user/add", ldap_server.AddUserLdapServer)
	r.GET("/ldapserver/user/search", ldap_server.SearchUserLdapServer)
	r.POST("/ldapserver/user/modify", ldap_server.ModifyUserLdapServer)
	r.POST("/ldapserver/user/delete/:cn", ldap_server.DeleteUserLdapServer)
	r.POST("/ldapserver/enterprise/init", ldap_server.InitEnterpriseLDAP)
	r.POST("/ldapserver/enterprise/clear", ldap_server.ClearEnterpriseLDAP)

	// Windows AD 路由
	r.POST("/winad/user/add", win_ad.AddUserWinAd)
	r.GET("/winad/user/search", win_ad.SearchUserWinAd)
	r.POST("/winad/user/modify", win_ad.ModifyUserWinAd)
	r.POST("/winad/user/delete/:cn", win_ad.DeleteUserWinAd)
	r.POST("/winad/enterprise/init", win_ad.InitEnterpriseWinAd)
	r.POST("/winad/enterprise/clear", win_ad.ClearEnterpriseWinAd)
	r.GET("/winad/enterprise/crawl", win_ad.CrawlEnterpriseWinAd)

	r.Run(":8080")
}
