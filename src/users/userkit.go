package users

import (
	"errors"
	"esnd/src/cry"
	"esnd/src/db"
	"regexp"
	"strings"
)

type User struct {
	Name string
	Md5  string
	Priv string
}

func Auth(name string, pw string) (*User, error) {

	reg, _ := regexp.Compile("^[0-9a-zA-Z_]{1,}$")
	if !reg.MatchString(name) {
		return nil, errors.New("invalid user name")
	}

	var u User
	u.Name = name
	u.Md5 = cry.MD5(pw)
	if name != "root" {
		if db.Count("SELECT count(*) FROM users WHERE name='"+name+"'") < 1 {
			return nil, errors.New("Auth Failed")
		}
		row := db.DB.QueryRow("SELECT mask,priv FROM users WHERE name='" + name + "'")
		var mask string
		err := row.Scan(&mask, &u.Priv)
		if err != nil {
			return nil, err
		}
		if mask != u.Md5 {
			return nil, errors.New("Auth Failed")
		}
		return &u, nil
	} else {
		if pw == db.Cfg.GetAnyway("root.mask", "changeMe") {
			u.Priv = "account pull push"
			return &u, nil
		} else {
			return nil, errors.New("root Auth Failed")
		}
	}
}

func (u *User) Can(priv string) bool {
	if priv == "pull" {
		return true
	}
	return strings.Contains(u.Priv, priv)
}
