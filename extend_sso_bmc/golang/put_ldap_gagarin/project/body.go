package main

import (
	pb "sso_center/gen/cluster-contract"
	"strings"
)

const (
	disableSSORequestGagarin = `{ "LDAP": { "ServiceEnabled": false } }`
	enableSSORequestGagarin  = `{
		"LDAP": {
		  "ServiceEnabled": true,
		  "Authentication": {
			"Username": "@SsoDn",
			"Password": "@SsoPassword",
			"AuthenticationType": "UsernameAndPassword"
		  },
		  "ServiceAddresses": [
			"ldap://@SsoAddress"
		  ],
		  "RemoteRoleMapping": [
			  {
				  "LocalRole": "Operator",
				  "RemoteGroup": "accounts"
			  }
		  ],
		  "LDAPService": {
			"SearchSettings": {
			  "BaseDistinguishedNames": [
				"@SsoRootDn"
			  ],
			  "GroupsAttribute": "gidNumber",
			  "UsernameAttribute": "cn"
			}
		  }
		}
	  }`
)

func createLDAPSetBody(ssoDn, ssoPassword string, state pb.SsoState) string {
	req := ""

	if state == pb.SsoState_SSO_STATE_ACTIVE {
		req = strings.Replace(enableSSORequestGagarin, "@SsoDn", ssoDn, 1)
		req = strings.Replace(req, "@SsoPassword", ssoPassword, 1)
		req = strings.Replace(req, "@SsoAddress", SSO_HOST+":389", 1)
		req = strings.Replace(req, "@SsoRootDn", BASE_DN, 1)
		req = strings.Replace(req, "@SsoGroupAttribute", GROUP_FORMAT, 1)
		req = strings.Replace(req, "@SsoUserNameAttribute", NAME_FORMAT, 1)
	} else {
		req = disableSSORequestGagarin
	}

	return req
}
