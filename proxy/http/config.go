package http

import (
	"strings"
	"net"
	
	"github.com/xtls/xray-core/common/protocol"
)

func (a *Account) Equals(another protocol.Account) bool {
	if account, ok := another.(*Account); ok {
		return a.Username == account.Username
	}
	return false
}

func (a *Account) AsAccount() (protocol.Account, error) {
	return a, nil
}

func (sc *ServerConfig) HasAccount(username, password string) (bool, string) {
	if sc.Accounts == nil {
		return false, ""
	}
	p :=  ""
	found := false
	if username != "" {
		p, found = sc.Accounts[username]
		if !found {
			return false, ""
		}
		if p == password {
			found = true
		}
	} else {
		for n, v := range sc.Accounts {
			ipPort := strings.Split(v, "|")
			for _, v2 := range ipPort {
				if v2 == password {
					ip := net.ParseIP(v2)
					if ip != nil {
						username = n
						found = true
					}
					break
				}
			}
		}
	}

	return found, username
}
