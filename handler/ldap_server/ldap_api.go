package ldap_server

import (
	"net/http"

	"github.com/coderZoe/ldap-test/handler"
	"github.com/gin-gonic/gin"
)

func AddUserLdapServer(ctx *gin.Context) {
	user := handler.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := addUserLdapServer(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "用户添加成功"})
}

func SearchUserLdapServer(ctx *gin.Context) {
	filter := ctx.Query("filter")
	result, err := searchUserLdapServer(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func ModifyUserLdapServer(ctx *gin.Context) {
	user := handler.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := modifyUserLdapServer(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "用户修改成功"})
}

func DeleteUserLdapServer(ctx *gin.Context) {
	userCn := ctx.Param("cn")
	if userCn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户CN不能为空"})
		return
	}

	err := deleteUserLdapServer(userCn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

func InitEnterpriseLDAP(ctx *gin.Context) {
	err := initEnterpriseLdapServer()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "企业LDAP结构初始化成功"})
}
func ClearEnterpriseLDAP(ctx *gin.Context) {
	err := clearEnterpriseLdapServer()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "企业LDAP结构清除成功"})
}
