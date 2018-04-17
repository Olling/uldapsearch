package main

import (
	"os"
	"fmt"
	"flag"
	"encoding/json"
	"bufio"
	"strings"
	"github.com/howeyc/gopass"
)

var (
	Basedn = flag.String("basedn","","(Required) The LDAP base DN")
	OutputAttribute = flag.String("outputattribute","dn","(Optional) Get a specific attribute from the LDAP")
	Filter = flag.String("filter","","(Optional) The LDAP search filter")
	FilterAttribute = flag.String("filterattribute","","(Optional) The LDAP search filter attribute")
	LdapHost = flag.String("ldaphost","","(Required) The LDAP host")
	LdapPort = flag.Int("ldapport",0,"(Required) The LDAP host port")
	Tls = flag.Bool("tls", true,"(Optional) Defines if the connection should be encrypted")
	TlsSkipCertVerification = flag.Bool("skipcertverification", false,"(Optional) Defines if the server certificate should verified")
	
	conf Configuration
)

type Configuration struct {
	Basedn string
	FilterAttribute string
	LdapHost string
	LdapPort int
	Tls bool
	TlsSkipCertVerification bool
}

func checkFlags() {
	if *Filter != "" && *FilterAttribute != "" {
		fmt.Println("Ignoring FilterAttribute as Filter is defined")
	}

    	if len(flag.Args()) != 1 && *Filter == "" {
        	fmt.Println("Only one nonspecific argument is allowed")
		os.Exit(1)
    	}
}

func loadConfiguration(path string) {
        file,_ := os.Open(path)
        decoder := json.NewDecoder(file)
        decoder.Decode(&conf)
}

func checkConfig() {
	//Check that the FilterAttribute or Filter is defined
	if conf.FilterAttribute == "" && *Filter == "" {
		conf.FilterAttribute = "sAMAccountName"
	}

	//Fill out Filter
	if *Filter == "" {
		*Filter = "(&(" + conf.FilterAttribute + "=" + flag.Args()[0] + "))"
	}

	//Check that LdapHost is defined
	if conf.LdapHost == "" && *LdapHost == "" {
		fmt.Println("No LDAP host (LdapHost) was provided")
		os.Exit(1)
	}

	//Overwrite LdapHost
	if *LdapHost != "" {
		conf.LdapHost = *LdapHost
	}
	
	//Check that LdapPort is defined
	if conf.LdapPort == 0 && *LdapPort == 0 {
		fmt.Println("No LDAP port (LdapPort) was provided")
		os.Exit(1)
	}

	//Overwrite LdapPort
	if *LdapPort != 0 {
		conf.LdapPort = *LdapPort
	}

	for _,arg := range os.Args[1:] {
		if arg == "-tls" || strings.HasPrefix(arg,"-tls=") {
			conf.Tls = *Tls
		}

		if arg == "-skipcertverification" || strings.HasPrefix(arg,"-skipcertverification=") {
			conf.TlsSkipCertVerification = *TlsSkipCertVerification
		}
	}
}

func promptUser () (string,string){
	fmt.Println("Please enter your LDAP information")

	//Username
	fmt.Printf("Username: ")
	reader := bufio.NewReader(os.Stdin)
	binduser,_ := reader.ReadString('\n')
	binduser = strings.TrimSpace(binduser)

	if !strings.Contains(binduser,conf.LdapHost) {
		binduser = binduser + "@" + conf.LdapHost
	}

	//Password
	fmt.Printf("Password: ")
	bindpass,_ := gopass.GetPasswdMasked()
	return binduser, string(bindpass)
}
