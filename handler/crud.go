package handler

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

func DoAddUser(user User, conn *ldap.Conn, baseDN string) error {
	//step1. 创建一个ou
	ouDN := "ou=users," + baseDN
	ou := ldap.NewAddRequest(ouDN, []ldap.Control{})
	ou.Attribute("objectClass", []string{"organizationalUnit"})

	err := conn.Add(ou)
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

func DoSearchUser(filter string, conn *ldap.Conn, baseDN string) (map[string][]string, error) {
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

func DoDeleteUser(userCn string, conn *ldap.Conn, baseDN string) error {
	userDN := "cn=" + userCn + ",ou=users," + baseDN
	ldapDelRequest := ldap.NewDelRequest(userDN, []ldap.Control{})
	return conn.Del(ldapDelRequest)
}

func DoModifyUser(User User, conn *ldap.Conn, baseDN string) error {
	userDN := "cn=" + User.CN + ",ou=users," + baseDN
	modifyRequest := ldap.NewModifyRequest(userDN, []ldap.Control{})
	if User.Mail != "" {
		modifyRequest.Replace("mail", []string{User.Mail})
	}
	if User.TelephoneNumber != "" {
		modifyRequest.Replace("telephoneNumber", []string{User.TelephoneNumber})
	}
	return conn.Modify(modifyRequest)
}
