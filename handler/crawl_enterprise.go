package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
)

func CrawlEnterprise(ctx *gin.Context) {
	entries, err := crawlEnterprise()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, entries)
}

// 爬取企业LDAP结构
func crawlEnterprise() ([]*ldap.Entry, error) {
	conn, err := connectAndBindLdap()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"*", "+"}, // 关键：获取所有用户属性和操作属性
		nil,
	)

	// 2. 执行搜索
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("全量搜索DN '%s' 失败: %w", baseDN, err)
	}

	fmt.Printf("✅ 抓取完成！共找到 %d 个条目。\n", len(sr.Entries))
	return sr.Entries, nil
}
