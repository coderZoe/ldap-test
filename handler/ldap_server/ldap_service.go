package ldap_server

import (
	"fmt"
	"log"

	"github.com/coderZoe/ldap-test/handler"
	"github.com/go-ldap/ldap/v3"
)

const (
	ldapHost      = "192.168.31.166"
	ldapPort      = 389
	adminDN       = "cn=admin,dc=example,dc=org"
	adminPassword = "adminpassword"
	baseDN        = "dc=example,dc=org"
)

func connectAndBindLdap() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", ldapHost, ldapPort))
	if err != nil {
		log.Printf("无法连接到LDAP服务器 ldap://%s:%d,失败原因 %v", ldapHost, ldapPort, err)
		return nil, err
	}

	err = conn.Bind(adminDN, adminPassword)
	if err != nil {
		log.Printf("无法绑定到LDAP服务器 ldap://%s:%d,失败原因 %v", ldapHost, ldapPort, err)
		return nil, err
	}
	fmt.Printf("连接LDAP服务器成功 ldap://%s:%d\n", ldapHost, ldapPort)
	return conn, nil
}

func addUserLdapServer(user handler.User) error {
	conn, err := connectAndBindLdap()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoAddUser(user, conn, baseDN)
}

func searchUserLdapServer(filter string) (map[string][]string, error) {
	conn, err := connectAndBindLdap()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return handler.DoSearchUser(filter, conn, baseDN)
}

func deleteUserLdapServer(userCn string) error {
	conn, err := connectAndBindLdap()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoDeleteUser(userCn, conn, baseDN)
}

func modifyUserLdapServer(user handler.User) error {
	conn, err := connectAndBindLdap()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoModifyUser(user, conn, baseDN)
}

func initEnterpriseLdapServer() error {
	conn, err := connectAndBindLdap()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoInitEnterprise(conn, baseDN)
}

func clearEnterpriseLdapServer() error {
	conn, err := connectAndBindLdap()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoClearEnterprise(conn, baseDN)
}

func crawlEnterprise() ([]*ldap.Entry, error) {
	conn, err := connectAndBindLdap()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	entries, err := handler.DoCrawlEnterprise(conn, baseDN)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
