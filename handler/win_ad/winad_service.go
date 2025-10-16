package win_ad

import (
	"fmt"
	"log"

	"github.com/coderZoe/ldap-test/handler"
	"github.com/go-ldap/ldap/v3"
)

//用于win server ad的LDAP操作

const (
	winAdHost          = "192.168.31.73"
	winAdPort          = 636
	winAdAdminDN       = "Administrator"
	winAdAdminPassword = "*********"
	winAdBaseDN        = "CN=Users,DC=example,DC=org"
)

func connectWinAd() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", winAdHost, winAdPort))
	if err != nil {
		log.Printf("无法连接到LDAP服务器 ldap://%s:%d,失败原因 %v", winAdHost, winAdPort, err)
		return nil, err
	}

	err = conn.Bind(winAdAdminDN, winAdAdminPassword)
	if err != nil {
		log.Printf("无法绑定到LDAP服务器 ldap://%s:%d,失败原因 %v", winAdHost, winAdPort, err)
		return nil, err
	}
	fmt.Printf("连接LDAP服务器成功 ldap://%s:%d\n", winAdHost, winAdPort)
	return conn, nil
}

func addUserWinAd(user handler.User) error {
	conn, err := connectWinAd()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoAddUser(user, conn, winAdBaseDN)
}

func searchUserWinAd(filter string) (map[string][]string, error) {
	conn, err := connectWinAd()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return handler.DoSearchUser(filter, conn, winAdBaseDN)
}

func deleteUserWinAd(userCn string) error {
	conn, err := connectWinAd()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoDeleteUser(userCn, conn, winAdBaseDN)
}

func modifyUserWinAd(user handler.User) error {
	conn, err := connectWinAd()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoModifyUser(user, conn, winAdBaseDN)
}

func initEnterpriseWinAd() error {
	conn, err := connectWinAd()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoInitEnterprise(conn, winAdBaseDN)
}

func clearEnterpriseWinAd() error {
	conn, err := connectWinAd()
	if err != nil {
		return err
	}
	defer conn.Close()

	return handler.DoClearEnterprise(conn, winAdBaseDN)
}

func crawlEnterpriseWinAd() ([]*ldap.Entry, error) {
	conn, err := connectWinAd()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return handler.DoCrawlEnterprise(conn, winAdBaseDN)
}
