package oidc

import (
	"authmantle-sso/data"
	"context"
)

type ObjectWrapper interface{}

var oidcSettings map[string]*ObjectWrapper

func initOIDCSettings() {
	oidcSettings = make(map[string]*ObjectWrapper)
}

func GetOIDCSetting(key string) *ObjectWrapper {
	return oidcSettings[key]
}

func Reload() {
	initOIDCSettings()
	//settings := GetOIDCSetting("oidc")
}

func BootStrapSettings(db *data.DatabaseHandler) {
	initOIDCSettings()
	connection, err := db.Acquire(context.Background())
	defer func() {
		connection.Release()
	}()
	if err != nil {
		return
	}
	//settings := GetOIDCSetting("oidc")
}
