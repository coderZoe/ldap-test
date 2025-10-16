package handler

import (
	"fmt"
	"slices"

	"github.com/go-ldap/ldap/v3"
)

// InitEnterpriseLDAP 初始化企业 LDAP 结构
func DoInitEnterprise(conn *ldap.Conn, baseDN string) error {
	// 创建组织单元
	if err := createOUs(conn, baseDN); err != nil {
		return err
	}

	// 创建部门
	if err := createDepartments(conn, baseDN); err != nil {
		return err
	}

	// 创建员工
	if err := createEmployees(conn, baseDN); err != nil {
		return err
	}

	// 创建资源
	if err := createResources(conn, baseDN); err != nil {
		return err
	}

	// 创建认证组
	if err := createAuthGroups(conn, baseDN); err != nil {
		return err
	}

	fmt.Println("企业 LDAP 结构初始化完成!")
	return nil
}

func createOUs(conn *ldap.Conn, baseDN string) error {
	ous := []struct {
		dn string
		ou string
	}{
		{"ou=Departments," + baseDN, "Departments"},
		{"ou=Resources," + baseDN, "Resources"},
		{"ou=Auth," + baseDN, "Auth"},
	}

	for _, item := range ous {
		req := ldap.NewAddRequest(item.dn, nil)
		req.Attribute("objectClass", []string{"top", "organizationalUnit"})
		req.Attribute("ou", []string{item.ou})
		if err := conn.Add(req); err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				fmt.Printf("⏭️  跳过已存在的OU: %s\n", item.dn)
				continue
			}
			fmt.Printf("OU %s 创建失败: %v\n", item.dn, err)
			return err
		} else {
			fmt.Printf("成功创建OU: %s\n", item.dn)
		}
	}
	return nil
}

func createDepartments(conn *ldap.Conn, baseDN string) error {
	departments := []struct {
		dn string
		ou string
	}{
		{"ou=Marketing,ou=Departments," + baseDN, "Marketing"},
		{"ou=Engineering,ou=Departments," + baseDN, "Engineering"},
		{"ou=HR,ou=Departments," + baseDN, "HR"},
		{"ou=Finance,ou=Departments," + baseDN, "Finance"},
		{"ou=Printers,ou=Resources," + baseDN, "Printers"},
		{"ou=Rooms,ou=Resources," + baseDN, "Rooms"},
		{"ou=Services,ou=Resources," + baseDN, "Services"},
	}

	for _, dept := range departments {
		req := ldap.NewAddRequest(dept.dn, nil)
		req.Attribute("objectClass", []string{"top", "organizationalUnit"})
		req.Attribute("ou", []string{dept.ou})

		if err := conn.Add(req); err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				fmt.Printf("⏭️  跳过已存在的部门: %s\n", dept.dn)
				continue
			}
			fmt.Printf("部门 %s 创建失败: %v\n", dept.dn, err)
			return err
		} else {
			fmt.Printf("成功创建部门: %s\n", dept.dn)
		}
	}
	return nil
}

func createEmployees(conn *ldap.Conn, baseDN string) error {
	employees := []struct {
		dn        string
		uid       string
		cn        string
		sn        string
		givenName string
		mail      string
		title     string
		manager   string
	}{
		// Marketing
		{"uid=alice.wang,ou=Marketing,ou=Departments," + baseDN, "alice.wang", "Alice Wang", "Wang", "Alice", "alice.wang@acme.com", "Marketing Specialist", ""},
		{"uid=bob.li,ou=Marketing,ou=Departments," + baseDN, "bob.li", "Bob Li", "Li", "Bob", "bob.li@acme.com", "Marketing Manager", ""},
		// Engineering
		{"uid=carol.zhang,ou=Engineering,ou=Departments," + baseDN, "carol.zhang", "Carol Zhang", "Zhang", "Carol", "carol.zhang@acme.com", "Backend Engineer", "uid=dave.chen,ou=Engineering,ou=Departments," + baseDN},
		{"uid=dave.chen,ou=Engineering,ou=Departments," + baseDN, "dave.chen", "Dave Chen", "Chen", "Dave", "dave.chen@acme.com", "Engineering Manager", ""},
		// HR
		{"uid=emma.liu,ou=HR,ou=Departments," + baseDN, "emma.liu", "Emma Liu", "Liu", "Emma", "emma.liu@acme.com", "HR Generalist", "uid=frank.huang,ou=HR,ou=Departments," + baseDN},
		{"uid=frank.huang,ou=HR,ou=Departments," + baseDN, "frank.huang", "Frank Huang", "Huang", "Frank", "frank.huang@acme.com", "HR Manager", ""},
		// Finance
		{"uid=grace.zhao,ou=Finance,ou=Departments," + baseDN, "grace.zhao", "Grace Zhao", "Zhao", "Grace", "grace.zhao@acme.com", "Accountant", "uid=henry.gu,ou=Finance,ou=Departments," + baseDN},
		{"uid=henry.gu,ou=Finance,ou=Departments," + baseDN, "henry.gu", "Henry Gu", "Gu", "Henry", "henry.gu@acme.com", "Finance Manager", ""},
	}

	for _, emp := range employees {
		req := ldap.NewAddRequest(emp.dn, nil)
		req.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "inetOrgPerson"})
		req.Attribute("uid", []string{emp.uid})
		req.Attribute("cn", []string{emp.cn})
		req.Attribute("sn", []string{emp.sn})
		req.Attribute("givenName", []string{emp.givenName})
		req.Attribute("mail", []string{emp.mail})
		req.Attribute("title", []string{emp.title})
		if emp.manager != "" {
			req.Attribute("manager", []string{emp.manager})
		}

		if err := conn.Add(req); err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				fmt.Printf("⏭️  跳过已存在的员工: %s\n", emp.dn)
				continue
			}
			fmt.Printf("创建员工 %s 失败: %v\n", emp.dn, err)
			return err
		} else {
			fmt.Printf("成功创建员工: %s\n", emp.cn)
		}
	}
	return nil
}

func createResources(conn *ldap.Conn, baseDN string) error {
	// 创建打印机
	printers := []struct {
		dn           string
		cn           string
		ipHostNumber string
		location     string
		description  string
	}{
		{"cn=printer-hq-1,ou=Printers,ou=Resources," + baseDN, "printer-hq-1", "10.0.10.21", "HQ-3F", "HQ三层东侧黑白激光打印机"},
		{"cn=printer-hq-2,ou=Printers,ou=Resources," + baseDN, "printer-hq-2", "10.0.10.22", "HQ-3F", "HQ三层彩色打印机"},
	}

	for _, printer := range printers {
		req := ldap.NewAddRequest(printer.dn, nil)
		req.Attribute("objectClass", []string{"top", "device", "ipHost"})
		req.Attribute("cn", []string{printer.cn})
		req.Attribute("ipHostNumber", []string{printer.ipHostNumber})
		req.Attribute("l", []string{printer.location})
		req.Attribute("description", []string{printer.description})

		if err := conn.Add(req); err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				fmt.Printf("⏭️  跳过已存在的打印机: %s\n", printer.dn)
				continue
			}
			fmt.Printf("创建打印机 %s 失败: %v\n", printer.dn, err)
			return err
		} else {
			fmt.Printf("成功创建打印机: %s\n", printer.cn)
		}
	}

	// 创建会议室
	rooms := []struct {
		dn          string
		cn          string
		roomNumber  string
		description string
	}{
		{"cn=room-Alpha,ou=Rooms,ou=Resources," + baseDN, "room-Alpha", "A301", "8人会议室,含电视/白板"},
		{"cn=room-Beta,ou=Rooms,ou=Resources," + baseDN, "room-Beta", "B302", "4人会议室,含电视"},
	}

	for _, room := range rooms {
		req := ldap.NewAddRequest(room.dn, nil)
		req.Attribute("objectClass", []string{"top", "room"})
		req.Attribute("cn", []string{room.cn})
		req.Attribute("roomNumber", []string{room.roomNumber})
		req.Attribute("description", []string{room.description})

		if err := conn.Add(req); err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				fmt.Printf("⏭️  跳过已存在的条目: %s\n", room.dn)
				continue
			}
			fmt.Printf("创建会议室 %s 失败: %v\n", room.dn, err)
			return err
		} else {
			fmt.Printf("成功创建会议室: %s\n", room.cn)
		}
	}

	// 创建应用服务
	req := ldap.NewAddRequest("cn=gitlab,ou=Services,ou=Resources,"+baseDN, nil)
	req.Attribute("objectClass", []string{"top", "applicationProcess"})
	req.Attribute("cn", []string{"gitlab"})
	req.Attribute("l", []string{"https://gitlab.acme.com/"})
	req.Attribute("description", []string{"Self-hosted GitLab"})

	if err := conn.Add(req); err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
			fmt.Printf("⏭️  跳过已存在的应用: %s\n", req.DN)
			return nil
		}
		fmt.Printf("创建应用 gitlab 失败: %v\n", err)
		return err
	} else {
		fmt.Println("成功创建应用: gitlab")
	}

	return nil
}

func createAuthGroups(conn *ldap.Conn, baseDN string) error {
	groups := []struct {
		dn      string
		cn      string
		members []string
	}{
		{
			"cn=gitlab-users,ou=Auth," + baseDN,
			"gitlab-users",
			[]string{
				"uid=carol.zhang,ou=Engineering,ou=Departments," + baseDN,
				"uid=dave.chen,ou=Engineering,ou=Departments," + baseDN,
			},
		},
		{
			"cn=jenkins-users,ou=Auth," + baseDN,
			"jenkins-users",
			[]string{
				"uid=carol.zhang,ou=Engineering,ou=Departments," + baseDN,
			},
		},
		{
			"cn=printer-hq-allowed,ou=Auth," + baseDN,
			"printer-hq-allowed",
			[]string{
				"uid=alice.wang,ou=Marketing,ou=Departments," + baseDN,
				"uid=grace.zhao,ou=Finance,ou=Departments," + baseDN,
			},
		},
	}

	for _, group := range groups {
		req := ldap.NewAddRequest(group.dn, nil)
		req.Attribute("objectClass", []string{"top", "groupOfNames"})
		req.Attribute("cn", []string{group.cn})
		req.Attribute("member", group.members)

		if err := conn.Add(req); err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				fmt.Printf("⏭️  跳过已存在的组: %s\n", group.dn)
				continue
			}
			fmt.Printf("创建组 %s 失败: %v\n", group.dn, err)
			return err
		} else {
			fmt.Printf("成功创建组: %s\n", group.cn)
		}
	}

	return nil
}
func DoClearEnterprise(conn *ldap.Conn, baseDN string) error {
	dnList, err := doGetAllSubDN(conn, baseDN)
	//逆序，从叶子节点删除
	slices.Reverse(dnList)
	if err != nil {
		return fmt.Errorf("无法获取所有子DN: %w", err)
	}
	for _, dn := range dnList {
		err := conn.Del(ldap.NewDelRequest(dn, []ldap.Control{}))
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			fmt.Printf("⏭️  跳过已不存在的条目: %s\n", dn)
			continue
		}
		if err != nil {
			return fmt.Errorf("无法删除DN: %s %w", dn, err)
		}
	}

	return nil
}

func doGetAllSubDN(conn *ldap.Conn, baseDN string) ([]string, error) {
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"dn"}, // 关键优化：只请求DN
		nil,
	)
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("无法搜索LDAP: %w", err)
	}

	dnList := []string{}
	for _, entry := range searchResult.Entries {
		dn := entry.DN
		dnList = append(dnList, dn)
	}
	return dnList, nil
}

// 爬取企业LDAP结构
func DoCrawlEnterprise(conn *ldap.Conn, baseDN string) ([]*ldap.Entry, error) {
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
