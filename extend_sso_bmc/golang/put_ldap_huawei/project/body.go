package main

import (
	pb "sso_center/gen/cluster-contract"
	"strings"
)

// Вставляю креды админа потому что хуавей не может авторизоваться с кредами тачки
const (
	disableLDAPBody  = `{ "LdapServiceEnabled": false }`
	enableLDAPBody   = `{ "LdapServiceEnabled": true }`
	settingsLDAPBody = `{
		"LdapServerAddress": "@LDAPAddress",
		"LdapPort": @LDAPPort,
		"UserDomain": ",DC=bergen.ems",
		"BindDN": "cn=nerd,dc=bergen,dc=ems",
		"BindPassword": "0penBmc",
		"CertificateVerificationEnabled ": false,
		"CertificateVerificationLevel": "Demand",
		"LdapGroups": [
			{
				"GroupName": "accounts",
				"GroupRole": "Operator",
				"GroupFolder": "OU=groups",
				"GroupLoginRule": [],
				"GroupLoginInterface": [
					"Web",
					"SSH",
					"Redfish"
				]
			}
		]
	  }`
	loadCABody = `{"Type":"text","Content": "@Cert"}`
)

func createLDAPManageBody(state pb.SsoState) string {
	if state == pb.SsoState_SSO_STATE_ACTIVE {
		return enableLDAPBody
	} else {
		return disableLDAPBody
	}
}

func createLDAPSetCABody(ca string) string {
	return strings.Replace(loadCABody, "@Cert", ca, 1)
}

func createLDAPSettingsBody() string {
	req := strings.Replace(settingsLDAPBody, "@LDAPAddress", SSO_HOST, 1)
	req = strings.Replace(req, "@LDAPPort", SSO_PORT, 1)

	return req
}
