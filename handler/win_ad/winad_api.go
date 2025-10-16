package win_ad

import (
	"net/http"

	"github.com/coderZoe/ldap-test/handler"
	"github.com/gin-gonic/gin"
)

func AddUserWinAd(ctx *gin.Context) {
	user := handler.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := addUserWinAd(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func SearchUserWinAd(ctx *gin.Context) {
	filter := ctx.Query("filter")
	result, err := searchUserWinAd(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func ModifyUserWinAd(ctx *gin.Context) {
	user := handler.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := modifyUserWinAd(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func DeleteUserWinAd(ctx *gin.Context) {
	userCn := ctx.Param("cn")
	if userCn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户CN不能为空"})
		return
	}
	err := deleteUserWinAd(userCn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func InitEnterpriseWinAd(ctx *gin.Context) {
	err := initEnterpriseWinAd()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "企业LDAP结构初始化成功"})
}

func ClearEnterpriseWinAd(ctx *gin.Context) {
	err := clearEnterpriseWinAd()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "企业LDAP结构清除成功"})
}

func CrawlEnterpriseWinAd(ctx *gin.Context) {
	entries, err := crawlEnterpriseWinAd()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, entries)
}
