package win_ad

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/coderZoe/ldap-test/handler"
	"github.com/go-ldap/ldap/v3"
)

//用于win server ad的LDAP操作

const (
	winAdHost    = "DC01.example.org"
	winAdPort    = 636
	winAdAdminDN = "cn=Administrator,cn=Users,dc=example,dc=org"

	winAdAdminPassword = "********"
	winAdBaseDN        = "dc=example,dc=org"

	caPemPath = "ldaps.pem"
)

func mustRootPool() *x509.CertPool {
	pem, err := os.ReadFile(caPemPath)
	if err != nil {
		log.Fatalf("读取 CA 失败(%s): %v", caPemPath, err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(pem) {
		log.Fatalf("加载 CA 失败: 不是有效 PEM")
	}
	return cp
}
func connectWinAd() (*ldap.Conn, error) {
	cp := mustRootPool()
	tlsCfg := &tls.Config{
		ServerName: winAdHost, // 与证书 SAN 匹配
		RootCAs:    cp,
		MinVersion: tls.VersionTLS12,
	}

	url := fmt.Sprintf("ldaps://%s:%d", winAdHost, winAdPort)

	conn, err := ldap.DialURL(
		url,
		ldap.DialWithTLSConfig(tlsCfg),
	)
	if err != nil {
		return nil, fmt.Errorf("LDAPS 连接失败: %w", err)
	}

	if err := conn.Bind(winAdAdminDN, winAdAdminPassword); err != nil {
		conn.Close()
		return nil, fmt.Errorf("bind win ad 失败: %w", err)
	}
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
