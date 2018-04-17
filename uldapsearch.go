package main

import (
	"os"
	"fmt"
	"flag"
	"strings"
	"crypto/tls"
	"gopkg.in/ldap.v2"
)

func ldapSearch(binduser string, bindpass string) {
	var err error
	var l *ldap.Conn
	if conf.Tls {
	       	l, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", conf.LdapHost, conf.LdapPort),&tls.Config{InsecureSkipVerify: conf.TlsSkipCertVerification, ServerName: conf.LdapHost})
	} else {
        	l, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", conf.LdapHost, conf.LdapPort))
	}

        if err != nil {
                fmt.Println(err)
		os.Exit(1)
        }

        defer l.Close()

        err = l.Bind(binduser, bindpass)
        if err != nil {
                fmt.Println(err)
		os.Exit(1)
        }

        searchRequest := ldap.NewSearchRequest(
                conf.Basedn,
                ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
                *Filter,
                strings.Split(*OutputAttribute,","),
                nil,
        )

        sr, err := l.Search(searchRequest)
        if err != nil {
                fmt.Println(err)
		os.Exit(1)
        }

	fmt.Println("Result:")
        for _, entry := range sr.Entries {
		result,_ := toJson(toJsonEntry(*entry))
		fmt.Println(result)
        }
}

func main () {
	//Set default Tls value
	conf.Tls = true

	//Loads default configuration from the harddrive
	loadConfiguration("/etc/uldapsearch.conf")

	//Gets the flags set by the user
	flag.Parse()
	checkFlags()

	//Check the configuration
	checkConfig()

	binduser,bindpass := promptUser()
        ldapSearch(binduser, string(bindpass))
}
