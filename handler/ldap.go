package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
)

const (
	ldapHost      = "192.168.31.166"
	ldapPort      = 389
	adminDN       = "cn=admin,dc=example,dc=org"
	adminPassword = "adminpassword"
	baseDN        = "dc=example,dc=org"
)

type User struct {
	CN              string `json:"cn"`
	SN              string `json:"sn"`
	Mail            string `json:"mail"`
	TelephoneNumber string `json:"telephoneNumber"`
}

func AddUser(ctx *gin.Context) {
	user := User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := doAddUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "用户添加成功"})
}

func SearchUser(ctx *gin.Context) {
	filter := ctx.Query("filter")
	result, err := doSearchUser(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func ModifyUser(ctx *gin.Context) {
	user := User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := doModifyUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "用户修改成功"})
}

func DeleteUser(ctx *gin.Context) {
	userCn := ctx.Param("cn")
	if userCn == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户CN不能为空"})
		return
	}

	err := doDeleteUser(userCn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

func doAddUser(user User) error {
	//step1. 创建一个ou
	ouDN := "ou=users," + baseDN
	ou := ldap.NewAddRequest(ouDN, []ldap.Control{})
	ou.Attribute("objectClass", []string{"organizationalUnit"})
	conn, err := connectAndBindLdap()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Add(ou)
	if err == nil {
		fmt.Println("成功创建ou users")
	} else {
		fmt.Printf("创建ou users失败，可能是ou已经存在 %v\n", err)
	}
	//正式添加用户
	userReq := ldap.NewAddRequest("cn="+user.CN+","+ouDN, []ldap.Control{})
	userReq.Attribute("objectClass", []string{"inetOrgPerson"})
	userReq.Attribute("cn", []string{user.CN})
	userReq.Attribute("sn", []string{user.SN})
	userReq.Attribute("mail", []string{user.Mail})
	userReq.Attribute("telephoneNumber", []string{user.TelephoneNumber})
	return conn.Add(userReq)
}

func doSearchUser(filter string) (map[string][]string, error) {
	conn, err := connectAndBindLdap()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	searchRequest := ldap.NewSearchRequest(
		baseDN,                 // 从基础DN开始搜索
		ldap.ScopeWholeSubtree, // 搜索范围：整个子树
		ldap.NeverDerefAliases, // 不处理别名
		0,                      // 大小限制 (0表示不限制)
		0,                      // 时间限制 (0表示不限制)
		false,                  // TypesOnly: false表示返回属性和值
		filter,                 // 过滤器
		[]string{"dn", "cn", "sn", "mail", "telephoneNumber"}, // 需要返回的属性列表
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]string)
	for _, entry := range searchResult.Entries {
		for _, attribute := range entry.Attributes {
			result[attribute.Name] = attribute.Values
		}
	}
	return result, nil
}

func doDeleteUser(userCn string) error {
	userDN := "cn=" + userCn + ",ou=users," + baseDN
	ldapDelRequest := ldap.NewDelRequest(userDN, []ldap.Control{})
	conn, err := connectAndBindLdap()
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Del(ldapDelRequest)
}

func doModifyUser(User User) error {
	userDN := "cn=" + User.CN + ",ou=users," + baseDN
	modifyRequest := ldap.NewModifyRequest(userDN, []ldap.Control{})
	if User.Mail != "" {
		modifyRequest.Replace("mail", []string{User.Mail})
	}
	if User.TelephoneNumber != "" {
		modifyRequest.Replace("telephoneNumber", []string{User.TelephoneNumber})
	}
	conn, err := connectAndBindLdap()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Modify(modifyRequest)
}

func connectAndBindLdap() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", ldapHost, ldapPort))
	if err != nil {
		log.Fatalf("无法连接到LDAP服务器 ldap://%s:%d,失败原因 %v", ldapHost, ldapPort, err)
		return nil, err
	}

	err = conn.Bind(adminDN, adminPassword)
	if err != nil {
		log.Fatalf("无法绑定到LDAP服务器 ldap://%s:%d,失败原因 %v", ldapHost, ldapPort, err)
		return nil, err
	}
	fmt.Printf("连接LDAP服务器成功 ldap://%s:%d\n", ldapHost, ldapPort)
	return conn, nil
}
