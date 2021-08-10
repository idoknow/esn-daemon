package users

import (
	"errors"
	"esnd/src/db"
	"strings"
)

type User struct {
	Name string
	Md5  string
	Priv string
}

func Auth(name string, md5 string) (*User, error) {
	var u User
	u.Name = name
	u.Md5 = md5
	if name != "root" {
		row := db.DB.QueryRow("SELECT mask,priv FROM users WHERE name='" + name + "'")
		var mask string
		err := row.Scan(&mask, &u.Priv)
		if err != nil {
			return nil, err
		}
		if mask != md5 {
			return nil, errors.New("Auth Failed")
		}
		return &u, nil
	} else {
		if md5 == db.Cfg.GetAnyway("root.mask", "changeMe") {
			u.Priv = "account pull push"
			return &u, nil
		} else {
			return nil, errors.New("Auth failed")
		}
	}
}

func (u *User) Can(priv string) bool {
	if priv == "pull" {
		return true
	}
	return strings.Contains(u.Priv, priv)
}
