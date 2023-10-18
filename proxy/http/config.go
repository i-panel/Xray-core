package http

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/xtls/xray-core/common/protocol"
)

// MemoryAccount is an account type converted from Account.
type MemoryAccount struct {
	Username string
	Password string
}

func hashString(username , pass string) string {
	hash := ""
	val := strings.TrimSpace(username) + strings.TrimSpace(pass)
	md5 := md5.Sum([]byte(strings.ToLower(val)))
	hash = fmt.Sprintf("%x", md5)
	return hash
}

// Equals implements protocol.Account.Equals().
func (a *MemoryAccount) Equals(another protocol.Account) bool {
	if account, ok := another.(*MemoryAccount); ok {
		return a.Password == account.Password
	}
	return false
}

// AsAccount implements protocol.AsAccount.
func (a *Account) AsAccount() (protocol.Account, error) {
	password := a.GetPassword()
	username := a.GetUsername()
	return &MemoryAccount{
		Password: password,
		Username: username,
	}, nil
}
