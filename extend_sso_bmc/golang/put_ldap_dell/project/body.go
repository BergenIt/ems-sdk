package main

import (
	pb "sso_center/gen/cluster-contract"
	"strings"
)

const (
	disableLDAP = `{"LDAP":{"ServiceEnabled":false}}`
	enableLDAP  = `{
		"LDAP": {
			"LDAPService": {
				"SearchSettings": {
					"BaseDistinguishedNames": [
						"@BaseDN"
					],
					"GroupNameAttribute": "cn"
				}
			},
			"RemoteRoleMapping": [
				{
					"RemoteGroup": "cn=accounts,ou=groups,dc=bergen,dc=ems",
                	"LocalRole": "Operator"
				}
			],
			"ServiceAddresses": [
				"@SsoHost"
			],
			"ServiceEnabled": true
		}
	}`
	setLDAPAttrs = `{
		"Attributes": {
			"LDAP.1.Port": @SsoPort,
			"LDAP.1.CertValidationEnable": "Disabled",
        	"LDAP.1.BindDN": "@BindDN",
        	"LDAP.1.SearchFilter": "objectclass=*",
			"LDAP.1.BindPassword": "@BindPassword"
		}
	}`
	setCA = `{
		"CertificateType": "CA",
		"SSLCertificateFile": "@Cert"
	}`
)

func createLDAPManageBody(state pb.SsoState) string {
	req := ""

	if state == pb.SsoState_SSO_STATE_ACTIVE {
		req = strings.Replace(enableLDAP, "@BaseDN", BASE_DN, 1)
		req = strings.Replace(req, "@SsoHost", SSO_HOST, 1)
	} else {
		req = disableLDAP
	}

	return req
}

func createLDAPAttrsBodyDell(ssoDn, ssoPassword string) string {
	req := strings.Replace(setLDAPAttrs, "@SsoPort", SSO_PORT, 1)
	req = strings.Replace(req, "@BindDN", ssoDn, 1)
	req = strings.Replace(req, "@BindPassword", ssoPassword, 1)

	return req
}

func createLoadCABody(ca string) string {
	return strings.Replace(setCA, "@Cert", ca, 1)
}
