package session

import (
	"github.com/darabuchi/utils/cache"
	"testing"
)

func TestGen(t *testing.T) {
	err := cache.Connect("127.0.0.1:6379", 1, "")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	type User struct {
		Username string `json:"username,omitempty" gorm:"column:username;not null;index:idx_admin_username,unique"`
		Password string `json:"password,omitempty" gorm:"column:password;not null"`
	}

	sid, err := GenSession(&User{
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(sid)

	data, err := GetSession(sid)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(data)

	var a User
	err = GetSessionJson(sid, &a)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(a)
}
